package domain

type MoneySetting struct {
	ChangeReserveCount  ChangeReserveCount  `json:"changeReserveCount"`  //釣銭準備金枚数
	ChangeShortageCount ChangeShortageCount `json:"changeShortageCount"` //不足枚数
	ExcessChangeCount   ExcessChangeCount   `json:"excessChangeCount"`   //あふれ枚数
}

type ChangeReserveCount struct {
	LastRegistDate string `json:"lastRegistDate"` //最終登録日付
	LastRegistTime string `json:"lastRegistTime"` //最終登録時刻
	M10000Count    int    `json:"m10000Count"`    //10000円枚数
	M5000Count     int    `json:"m5000Count"`     //5000円枚数
	M2000Count     int    `json:"m2000Count"`     //2000円枚数
	M1000Count     int    `json:"m1000Count"`     //1000円枚数
	M500Count      int    `json:"m500Count"`      //500円枚数
	M100Count      int    `json:"m100Count"`      //100円枚数
	M50Count       int    `json:"m50Count"`       //50円枚数
	M10Count       int    `json:"m10Count"`       //10円枚数
	M5Count        int    `json:"m5Count"`        //5C円枚数
	M1Count        int    `json:"m1Count"`        //1C円枚数
	S500Count      int    `json:"s500Count"`      //500円枚数(サブ)
	S100Count      int    `json:"s100Count"`      //100円枚数(サブ)
	S50Count       int    `json:"s50Count"`       //5C円枚数(サブ)
	S10Count       int    `json:"s10Count"`       //10円枚数(サブ)
	S5Count        int    `json:"s5Count"`        //5C円枚数(サブ)
	S1Count        int    `json:"s1Count"`        //1C円枚数(サブ)
}

type ChangeShortageCount struct {
	LastRegistDate  string             `json:"lastRegistDate"` //最終登録日付
	LastRegistTime  string             `json:"lastRegistTime"` //最終登録時刻
	RegisterDataTbl [2]RegisterDataTbl `json:"registerDataTbl"`
}

type RegisterDataTbl struct {
	AlertLevel  int `json:"alertLevel"`  //アラートレベル
	M10000Count int `json:"m10000Count"` //10000円枚数
	M5000Count  int `json:"m5000Count"`  //5000円枚数
	M2000Count  int `json:"m2000Count"`  //2000円枚数
	M1000Count  int `json:"m1000Count"`  //1000円枚数
	M500Count   int `json:"m500Count"`   //500円枚数
	M100Count   int `json:"m100Count"`   //100円枚数
	M50Count    int `json:"m50Count"`    //50円枚数
	M10Count    int `json:"m10Count"`    //10円枚数
	M5Count     int `json:"m5Count"`     //5C円枚数
	M1Count     int `json:"m1Count"`     //1C円枚数
	S500Count   int `json:"s500Count"`   //500円枚数(サブ)
	S100Count   int `json:"s100Count"`   //100円枚数(サブ)
	S50Count    int `json:"s50Count"`    //5C円枚数(サブ)
	S10Count    int `json:"s10Count"`    //10円枚数(サブ)
	S5Count     int `json:"s5Count"`     //5C円枚数(サブ)
	S1Count     int `json:"s1Count"`     //1C円枚数(サブ)
}

type ExcessChangeCount struct {
	LastRegistDate    string               `json:"lastRegistDate"` //最終登録日付
	LastRegistTime    string               `json:"lastRegistTime"` //最終登録時刻
	ExRegisterDataTbl [2]ExRegisterDataTbl `json:"registerDataTbl"`
}

type ExRegisterDataTbl struct {
	AlertLevel  int `json:"alertLevel"`  //アラートレベル
	M10000Count int `json:"m10000Count"` //10000円枚数
	M5000Count  int `json:"m5000Count"`  //5000円枚数
	M2000Count  int `json:"m2000Count"`  //2000円枚数
	M1000Count  int `json:"m1000Count"`  //1000円枚数
	M500Count   int `json:"m500Count"`   //500円枚数
	M100Count   int `json:"m100Count"`   //100円枚数
	M50Count    int `json:"m50Count"`    //50円枚数
	M10Count    int `json:"m10Count"`    //10円枚数
	M5Count     int `json:"m5Count"`     //5C円枚数
	M1Count     int `json:"m1Count"`     //1C円枚数
	S500Count   int `json:"s500Count"`   //500円枚数(サブ)
	S100Count   int `json:"s100Count"`   //100円枚数(サブ)
	S50Count    int `json:"s50Count"`    //5C円枚数(サブ)
	S10Count    int `json:"s10Count"`    //10円枚数(サブ)
	S5Count     int `json:"s5Count"`     //5C円枚数(サブ)
	S1Count     int `json:"s1Count"`     //1C円枚数(サブ)
	BillOverBox int `json:"billOverBox"` //全紙幣
	CoinOverBox int `json:"coinOverBox"` //全硬貨
}
