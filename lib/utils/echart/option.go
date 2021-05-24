// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package echart

type Option struct {
	Grid struct {
		Bottom       string `json:"bottom"`
		ContainLabel bool   `json:"containLabel"`
		Left         string `json:"left"`
		Right        string `json:"right"`
	} `json:"grid"`
	Legend struct {
		Data []string `json:"data"`
	} `json:"legend"`
	Series []struct {
		AreaStyle struct{} `json:"areaStyle,omitempty"`
		Data      []int64  `json:"data,omitempty"`
		Emphasis  struct {
			Focus string `json:"focus,omitempty"`
		} `json:"emphasis,omitempty"`
		Label struct {
			Position string `json:"position,omitempty"`
			Show     bool   `json:"show,omitempty"`
		} `json:"label,omitempty"`
		Name  string `json:"name,omitempty"`
		Stack string `json:"stack,omitempty"`
		Type  string `json:"type,omitempty"`
	} `json:"series,omitempty"`
	Title struct {
		Text string `json:"text,omitempty"`
	} `json:"title,omitempty"`
	Toolbox struct {
		Feature struct {
			SaveAsImage struct{} `json:"saveAsImage,omitempty"`
		} `json:"feature,omitempty"`
	} `json:"toolbox,omitempty"`
	Tooltip struct {
		AxisPointer struct {
			Label struct {
				BackgroundColor string `json:"backgroundColor,omitempty"`
			} `json:"label,omitempty"`
			Type string `json:"type,omitempty"`
		} `json:"axisPointer,omitempty"`
		Trigger string `json:"trigger,omitempty"`
	} `json:"tooltip,omitempty"`
	XAxis []struct {
		BoundaryGap bool     `json:"boundaryGap,omitempty"`
		Data        []string `json:"data,omitempty"`
		Type        string   `json:"type,omitempty"`
	} `json:"xAxis,omitempty"`
	YAxis []struct {
		Type string `json:"type,omitempty"`
	} `json:"yAxis,omitempty"`
}

func NewEChartOption() (res *Option) {
	return &Option{}
}
