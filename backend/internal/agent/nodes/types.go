package nodes

// GraphData 图数据传递结构
type GraphData struct {
	// 输入数据
	Image    string   // base64编码的图片
	Text     string   // 文本内容
	Age      int      // 年龄
	Keywords []string // 关键词

	// 中间数据
	ObjectName     string // 识别的对象名称
	ObjectCategory string // 对象类别
	Intent         string // 识别的意图

	// 输出数据
	Cards      []interface{} // 生成的卡片
	TextResult string        // 文本生成结果
	ImageURL   string        // 生成的图片URL
}
