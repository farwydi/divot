package warcraftlogs

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/farwydi/divot/pkg/enums"
	"google.golang.org/protobuf/proto"
)

var rankingVersion = []byte("v10")

func NewRankingKey(
	sessionId int,
	className string,
	specName string,
	encounterId int,
	reportCode string,
	fightID uint32,
	characterName string,
	server string,
) []byte {
	return bytes.Join([][]byte{
		[]byte("ranking"),
		rankingVersion,
		[]byte(strconv.Itoa(sessionId)),
		[]byte(className),
		[]byte(specName),
		[]byte(strconv.Itoa(encounterId)),
		[]byte(reportCode),
		[]byte(strconv.Itoa(int(fightID))),
		[]byte(characterName + "-" + server),
	}, []byte("/"))
}

func ParseRankingKey(key []byte) (
	sessionId int,
	className string,
	specName string,
	encounterId int,
	reportCode string,
	fightID int,
	characterFullName string,
) {
	elements := bytes.Split(key, []byte("/"))
	if !bytes.Equal(elements[0], []byte("ranking")) {
		panic("invalid ranking key")
	}

	var err error

	switch {
	case bytes.Equal(elements[1], rankingVersion):
		sessionId, err = strconv.Atoi(string(elements[2]))
		if err != nil {
			panic("invalid sessionId: " + err.Error())
		}

		className = string(elements[3])

		specName = string(elements[4])

		encounterId, err = strconv.Atoi(string(elements[5]))
		if err != nil {
			panic("invalid encounterId: " + err.Error())
		}

		reportCode = string(elements[6])

		fightID, err = strconv.Atoi(string(elements[7]))
		if err != nil {
			panic("invalid fightID: " + err.Error())
		}

		characterFullName = string(elements[8])

		return
	default:
		panic("invalid ranking version")
	}
}

type WorldDataEncounterCharacterRankings struct {
	WorldData struct {
		Encounter struct {
			CharacterRankings *CharacterRankingPage `scalar:"true" graphql:"characterRankings(includeCombatantInfo: true, page: $page, specName: $spec, className: $class, serverRegion: $region, metric: dps)"`
		} `graphql:"encounter(id: $encounterId)"`
	}
}

func (t *WarcraftLogs) ScanWorldDataEncounterCharacterRankings(ctx context.Context, encounterId int, className, specName string) error {
	page := 1
	for {
		if page > 20 {
			// done
			return nil
		}

		fmt.Printf("getting page %d\n", page)

		var q WorldDataEncounterCharacterRankings
		variables := map[string]interface{}{
			"encounterId": encounterId,
			"page":        page,
			"class":       enums.Fix.String(className),
			"spec":        enums.Fix.String(specName),
			"region":      "EU",
		}
		err := t.client.Query(ctx, &q, variables)
		if err != nil {
			return fmt.Errorf("fail query wow logs: %v with %v", err, variables)
		}

		err = t.saveRankings(encounterId, q.WorldData.Encounter.CharacterRankings.Rankings)
		if err != nil {
			return err
		}

		if !q.WorldData.Encounter.CharacterRankings.HasMorePages {
			break
		}

		page++
	}

	return nil
}

func (t *WarcraftLogs) saveRankings(encounterId int, logOnlyRankings []*Ranking) error {
	for _, ranking := range logOnlyRankings {
		if ranking.Name == "Anonymous" {
			if !t.saveAnonymous {
				continue
			}

			buf := make([]byte, 15)
			_, err := rand.Read(buf)
			if err != nil {
				panic(fmt.Errorf("crypto read fail: %v", err))
			}

			ranking.Server.Name = hex.EncodeToString(buf)
		}

		data, err := proto.Marshal(ranking)
		if err != nil {
			return err
		}

		key := NewRankingKey(
			1,
			ranking.Class,
			ranking.Spec,
			encounterId,
			ranking.Report.Code,
			ranking.Report.FightID,
			ranking.Name,
			ranking.Server.Name,
		)

		err = t.db.Write(key, data)
		if err != nil {
			return fmt.Errorf("write ranking: %v", err)
		}
	}

	return nil
}

func (t *WarcraftLogs) ReadWorldDataEncounterCharacterRankings(rankingProc func(encounterId int, ranking *Ranking) error) error {
	startTime := time.Now()
	defer func() {
		fmt.Printf("loaded %s\n", time.Since(startTime))
	}()

	prefix := bytes.Join([][]byte{[]byte("ranking"), rankingVersion, []byte("1")}, []byte("/"))

	return t.db.Scan(prefix, func(key, value []byte) error {
		_, _, _, encounterId, _, _, _ := ParseRankingKey(key)

		ranking := new(Ranking)
		err := proto.Unmarshal(value, ranking)
		if err != nil {
			return err
		}

		return rankingProc(encounterId, ranking)
	})
}
