package apieastmoney

import (
	"strings"

	"example.com/stocker-back/internal/infra"
)

type APIServiceEastmoney struct {
	logger infra.Logger
}

func NewAPIServiceEastmoney(logger infra.Logger) *APIServiceEastmoney {
	return &APIServiceEastmoney{
		logger: logger,
	}
}

// sliceStringByChar slice input string by startChar and endChar if they are valid.
func sliceStringByChar(input, startChar, endChar string) string {
	startIndex := strings.Index(input, startChar)
	if startIndex == -1 {
		return ""
	}

	endIndex := strings.LastIndex(input, endChar)
	if endIndex == -1 {
		return ""
	}

	return input[startIndex+1 : endIndex]
}

func SECTORS() []string {
	return []string{
		"食品饮料",
		"农牧饲渔",
		"玻璃玻纤",
		"仪器仪表",
		"石油行业",
		"医疗服务",
		"光伏设备",
		"软件开发",
		"航空机场",
		"电源设备",
		"农药兽药",
		"交运设备",
		"汽车零部件",
		"通用设备",
		"汽车服务",
		"包装材料",
		"装修装饰",
		"环保行业",
		"家电行业",
		"化学制品",
		"旅游酒店",
		"通信服务",
		"文化传媒",
		"房地产开发",
		"造纸印刷",
		"塑料制品",
		"电子化学品",
		"半导体",
		"化肥行业",
		"采掘行业",
		"珠宝首饰",
		"酿酒行业",
		"电力行业",
		"物流行业",
		"化学原料",
		"非金属材料",
		"美容护理",
		"电机",
		"消费电子",
		"光学光电子",
		"专业服务",
		"公用事业",
		"贸易行业",
		"保险",
		"贵金属",
		"电网设备",
		"船舶制造",
		"铁路公路",
		"游戏",
		"专用设备",
		"证券",
		"装修建材",
		"工程机械",
		"生物制品",
		"家用轻工",
		"纺织服装",
		"中药",
		"多元金融",
		"综合行业",
		"工程咨询服务",
		"医疗器械",
		"钢铁行业",
		"医药商业",
		"燃气",
		"煤炭行业",
		"汽车整车",
		"工程建设",
		"银行",
		"互联网服务",
		"航天航空",
		"化纤行业",
		"能源金属",
		"电池",
		"通信设备",
		"小金属",
		"水泥建材",
		"商业百货",
		"风电设备",
		"计算机设备",
		"化学制药",
		"电子元件",
		"航运港口",
		"有色金属",
		"橡胶制品",
		"教育",
	}
}
