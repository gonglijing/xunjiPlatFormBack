package model

// DictTypeOut 字典类型输出字段定义
type DictTypeOut struct {
	DictName string `json:"name"`
	Remark   string `json:"remark"`
}

// DictDataOut 字典数据输出字段定义
type DictDataOut struct {
	DictValue string `json:"key"`
	DictLabel string `json:"value"`
	IsDefault int    `json:"isDefault"`
	Remark    string `json:"remark"`
}

// GetDictInput 获取字典数据输入参数定义
type GetDictInput struct {
	DictType     string `json:"dictType"`
	DefaultValue string `json:"defaultValue"`
}

// GetDictOut 获取字典数据输出字段定义
type GetDictOut struct {
	Data   *DictTypeOut   `json:"info"`
	Values []*DictDataOut `json:"values"`
}
