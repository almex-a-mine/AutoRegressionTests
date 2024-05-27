package domain

// お金の価値
const (
	TEN_THOUSAND  = 10000
	FIVE_THOUSAND = 5000
	TWO_THOUSAND  = 2000
	ONE_THOUSAND  = 1000
	FIVE_HUNDRED  = 500
	ONE_HUNDRED   = 100
	FIFTY         = 50
	TEN           = 10
	FIVE          = 5
	ONE           = 1
)

// お金の価値配列
var Cash = [CASH_TYPE_SHITEI]int{
	TEN_THOUSAND,
	FIVE_THOUSAND,
	TWO_THOUSAND,
	ONE_THOUSAND,
	FIVE_HUNDRED,
	ONE_HUNDRED,
	FIFTY,
	TEN,
	FIVE,
	ONE}

//お金の価値配列拡張金種別枚数
/*var Extra_Cash = [EXTRA_CASH_TYPE_SHITEI]int{
TEN_THOUSAND,  //[0]:10,000円
FIVE_THOUSAND, //[1]:5,000円
TWO_THOUSAND,  //[2]:2,000円
ONE_THOUSAND,  //[3]:1,000円
FIVE_HUNDRED,  //[4]:500円
ONE_HUNDRED,   //[5]:100円
FIFTY,         //[6]:50円
TEN,           //[7]:10円
FIVE,          //[8]:5円
ONE,           //[9]:1円
FIVE_HUNDRED,  //[10]:500円予備
ONE_HUNDRED,   //[11]:100円予備
FIFTY,         //[12]:50円予備
TEN,           //[13]:10円予備
FIVE,          //[14]:5円予備
ONE,           //[15]:1円予備
TEN_THOUSAND,  //[16]:10,000円あふれ
FIVE_THOUSAND, //[17]:5,000円あふれ
TWO_THOUSAND,  //[18]:2,000円あふれ
ONE_THOUSAND,  //[19]:1,000円あふれ
FIVE_HUNDRED,  //[20]:500円あふれ
ONE_HUNDRED,   //[21]100円あふれ
FIFTY,         //[22]:50円あふれ
TEN,           //[23]:10円あふれ
FIVE,          //[24]:5円あふれ
ONE}           //[25]:1円あふれ*/

// tex_2_1用金種配列
var TexHelperCash = [CASH_TYPE_UI]int{
	TEN_THOUSAND,  //[0]:10,000円
	FIVE_THOUSAND, //[1]:5,000円
	TWO_THOUSAND,  //[2]:2,000円
	ONE_THOUSAND,  //[3]:1,000円
	FIVE_HUNDRED,  //[4]:500円
	ONE_HUNDRED,   //[5]:100円
	FIFTY,         //[6]:50円
	TEN,           //[7]:10円
	FIVE,          //[8]:5円
	ONE,           //[9]:1円
	FIVE_HUNDRED,  //[10]:500円予備
	ONE_HUNDRED,   //[11]:100円予備
	FIFTY,         //[12]:50円予備
	TEN,           //[13]:10円予備
	FIVE,          //[14]:5円予備
	ONE,           //[15]:1円予備
}

// １系金種のみの配列
var OneCashOnly = [5]int{TEN_THOUSAND, ONE_THOUSAND, ONE_HUNDRED, TEN, ONE}

// 1&5系金種のみの配列
var OnendFiveCashOnly = [9]int{
	TEN_THOUSAND,
	FIVE_THOUSAND,
	ONE_THOUSAND,
	FIVE_HUNDRED,
	ONE_HUNDRED,
	FIFTY,
	TEN,
	FIVE,
	ONE}

// 印刷制御用
// 精算機内合計用
var AllCashInMachine = [PRINT_CASH_DATA]int{
	TEN_THOUSAND,
	FIVE_THOUSAND,
	TWO_THOUSAND,
	ONE_THOUSAND,
	FIVE_HUNDRED,
	ONE_HUNDRED,
	FIFTY,
	TEN,
	FIVE,
	ONE,
	ONE_HUNDRED,
	TEN,
	ONE}

// 金庫
var Safe = [PRINT_CASH_DATA_SIX]int{
	FIVE_HUNDRED,
	ONE_HUNDRED,
	FIFTY,
	TEN,
	FIVE,
	ONE}

// 紙幣補充
var BillReplenishCash = [PRINT_CASH_DATA_SIXTEEN]int{
	0,
	0,
	0,
	0,
	TEN_THOUSAND,
	FIVE_THOUSAND,
	TWO_THOUSAND,
	ONE_THOUSAND,
	TEN_THOUSAND,
	FIVE_THOUSAND,
	TWO_THOUSAND,
	ONE_THOUSAND,
	0,
	0,
	0,
	0}

var AllCashInMachineTwentySix = [PRINT_CASH_DATA_EXCOUNT]int{
	TEN_THOUSAND,
	FIVE_THOUSAND,
	TWO_THOUSAND,
	ONE_THOUSAND,
	FIVE_HUNDRED, //通常金庫
	ONE_HUNDRED,
	FIFTY,
	TEN,
	FIVE,
	ONE,
	FIVE_HUNDRED, //予備金庫
	ONE_HUNDRED,
	FIFTY,
	TEN,
	FIVE,
	ONE,
	TEN_THOUSAND, //あふれ紙幣
	FIVE_THOUSAND,
	TWO_THOUSAND,
	ONE_THOUSAND,
	FIVE_HUNDRED, //硬貨入金庫枚数
	ONE_HUNDRED,
	FIFTY,
	TEN,
	FIVE,
	ONE}

var AllCashInMachineNine = [9]int{
	FIVE_HUNDRED, //通常金庫
	ONE_HUNDRED,
	FIFTY,
	TEN,
	FIVE,
	ONE,
	ONE_HUNDRED, //予備金庫
	TEN,
	ONE}

// 数の桁
const (
	TEN_THOUSAND_KETA = 10000
	ONE_THOUSAND_KETA = 1000
	ONE_HUNDRED_KETA  = 100
	TEN_KETA          = 10
	ONE_KETA          = 1
)

//数の桁配列
// var CashKeta = [5]int{TEN_THOUSAND_KETA, ONE_THOUSAND_KETA, ONE_HUNDRED_KETA, TEN_KETA, ONE_KETA}
