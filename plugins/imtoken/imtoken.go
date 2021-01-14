package imtoken

import (
	"fmt"
	ui "github.com/itering/subscan-plugin"
	"github.com/itering/subscan-plugin/router"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins/imtoken/dao"
	"github.com/itering/subscan/plugins/imtoken/http"
	"github.com/itering/subscan/plugins/imtoken/model"
	"github.com/itering/subscan/plugins/imtoken/service"
	"github.com/itering/subscan/util"
	"github.com/prometheus/common/log"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
)

var srv *service.Service

type Imtoken struct {
	d storage.Dao
}

func New() *Imtoken {
	return &Imtoken{}
}

func (a *Imtoken) InitDao(d storage.Dao) {
	srv = service.New(d)
	a.d = d
	a.Migrate()
}

func (a *Imtoken) InitHttp() []router.Http {
	return http.Router(srv)
}

func (a *Imtoken) ProcessExtrinsic(_ *storage.Block, _ *storage.Extrinsic, _ []storage.Event) error {
	return nil
}

func (a *Imtoken) ProcessEvent(block *storage.Block, event *storage.Event, fee decimal.Decimal) error {
	if event == nil {
		return nil
	}
	var paramEvent []storage.EventParam
	util.UnmarshalAny(&paramEvent, event.Params)

	switch fmt.Sprintf("%s-%s", strings.ToLower(event.ModuleId), strings.ToLower(event.EventId)) {
	case strings.ToLower("balances-transfer"):
		// 只存储 balances-transfer 的event
		if len(paramEvent) < 3 {
			log.Warn("parse balance transfer event param error", "param", paramEvent)
			return nil
		}
		transfer := &model.EventTransfer{
			EventIndex: strconv.Itoa(event.BlockNum) + "-" + strconv.Itoa(event.EventIdx),
			BlockNum:   uint64(event.BlockNum),
			BlockHash:  block.Hash,
			Timestamp:  uint64(block.BlockTimestamp),
			Sender:     util.ToString(paramEvent[0].Value),
			Receiver:   util.ToString(paramEvent[1].Value),
			Amount:     util.ToString(paramEvent[2].Value),
			// Nonce:              , // 需要通过 Extrinsic 查询
			Fee:           fee.String(),
			ExtrinsicHash: util.AddHex(event.ExtrinsicHash),
			ExtrinsicIdx:  event.ExtrinsicIdx,
			// CallModule:         "", // 需要通过 Extrinsic 查询
			// CallModuleFunction: "", // 需要通过 Extrinsic 查询
			ModuleId: event.ModuleId,
			EventId:  event.EventId,
			EventIdx: event.EventIdx,
		}
		return dao.NewTransfer(a.d, transfer)
	}

	return nil
}

func (a *Imtoken) Migrate() {
	_ = a.d.AutoMigration(&model.EventTransfer{})
	_ = a.d.AddUniqueIndex(&model.EventTransfer{}, "id", "id")
	_ = a.d.AddUniqueIndex(&model.EventTransfer{}, "event_index", "event_index")
	_ = a.d.AddIndex(&model.EventTransfer{}, "block_num", "block_num")
	_ = a.d.AddIndex(&model.EventTransfer{}, "sender", "sender")
	_ = a.d.AddIndex(&model.EventTransfer{}, "receiver", "receiver")
	_ = a.d.AddIndex(&model.EventTransfer{}, "extrinsic_hash", "extrinsic_hash")
}

// Plugins version
func (a *Imtoken) Version() string {
	return "0.1"
}

// Aims UI config
func (a *Imtoken) UiConf() *ui.UiConfig {
	return nil
}

// Subscribe Extrinsic with special module
func (a *Imtoken) SubscribeExtrinsic() []string {
	return nil
}

// Subscribe Events with special module
func (a *Imtoken) SubscribeEvent() []string {
	return []string{"balances"}
}
