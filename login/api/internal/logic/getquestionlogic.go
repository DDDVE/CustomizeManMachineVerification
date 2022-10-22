package logic

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zwh203080815/cmmvplat/call"
)

var randomSeed int64 = 0

const (
	QUESTION_TYPE        = 0
	FIRM_ID              = "testID"
	PRIVATE_KEY_PATH     = "../pkg/conf/eccprivate.pem"
	RANDOM_ANSWER_LENGTH = 6

	EXPIRE_LENGHT    = 70
	KEY_PREFIX_TOKEN = "cache:login:token:"

	TOKEN_SOURCE_STRING = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	COMMON_HANZI = "的一是了我不人在他有这个上们来到时大地为子中你说生国年着就那和要她出也得里后自以会家可下而过天去能对小多然于心学么之都好看起发当没成只如事把还用第样道想作种开美总从无情己面最女但现前些所同日手又行意动方期它头经长儿回位分爱老因很给名法间斯知世什两次使身者被高已亲其进此话常与活正感见明问力理尔点文几定本公特做外孩相西果走将月十实向声车全信重三机工物气每并别真打太新比才便夫再书部水像眼等体却加电主界门利海受听表德少克代员许稜先口由死安写性马光白或住难望教命花结乐色更拉东神记处让母父应直字场平报友关放至张认接告入笑内英军候民岁往何度山觉路带万男边风解叫任金快原吃妈变通师立象数四失满战远格士音轻目条呢病始达深完今提求清王化空业思切怎非找片罗钱紶吗语元喜曾离飞科言干流欢约各即指合反题必该论交终林请医晚制球决窢传画保读运及则房早院量苦火布品近坐产答星精视五连司巴奇管类未朋且婚台夜青北队久乎越观落尽形影红爸百令周吧识步希亚术留市半热送兴造谈容极随演收首根讲整式取照办强石古华諣拿计您装似足双妻尼转诉米称丽客南领节衣站黑刻统断福城故历惊脸选包紧争另建维绝树系伤示愿持千史谁准联妇纪基买志静阿诗独复痛消社算义竟确酒需单治卡幸兰念举仅钟怕共毛句息功官待究跟穿室易游程号居考突皮哪费倒价图具刚脑永歌响商礼细专黄块脚味灵改据般破引食仍存众注笔甚某沉血备习校默务土微娘须试怀料调广蜖苏显赛查密议底列富梦错座参八除跑亮假印设线温虽掉京初养香停际致阳纸李纳验助激够严证帝饭忘趣支春集丈木研班普导顿睡展跳获艺六波察群皇段急庭创区奥器谢弟店否害草排背止组州朝封睛板角况曲馆育忙质河续哥呼若推境遇雨标姐充围案伦护冷警贝著雪索剧啊船险烟依斗值帮汉慢佛肯闻唱沙局伯族低玩资屋击速顾泪洲团圣旁堂兵七露园牛哭旅街劳型烈姑陈莫鱼异抱宝权鲁简态级票怪寻杀律胜份汽右洋范床舞秘午登楼贵吸责例追较职属渐左录丝牙党继托赶章智冲叶胡吉卖坚喝肉遗救修松临藏担戏善卫药悲敢靠伊村戴词森耳差短祖云规窗散迷油旧适乡架恩投弹铁博雷府压超负勒杂醒洗采毫嘴毕九冰既状乱景席珍童顶派素脱农疑练野按犯拍征坏骨余承置臓彩灯巨琴免环姆暗换技翻束增忍餐洛塞缺忆判欧层付阵玛批岛项狗休懂武革良恶恋委拥娜妙探呀营退摇弄桌熟诺宣银势奖宫忽套康供优课鸟喊降夏困刘罪亡鞋健模败伴守挥鲜财孤枪禁恐伙杰迹妹藸遍盖副坦牌江顺秋萨菜划授归浪听凡预奶雄升碃编典袋莱含盛济蒙棋端腿招释介烧误"
)

type GetQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionLogic {
	return &GetQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetQuestionLogic) GetQuestion(req *types.QuestionReq) (resp *types.QuestionReply, err error) {

	//调用output接口
	data, err := call.GetQuestionByPath(QUESTION_TYPE, FIRM_ID, PRIVATE_KEY_PATH)
	if err != nil {
		return nil, err
	}

	//打乱答案顺序
	randomAnswer := RandomAnswer(data.Answer, data.DisturbAnswer)

	//生成token
	stamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	b := sha256.Sum256(append(randomAnswer, []byte(stamp)...))
	token := [len(b)]byte{}
	for i := 0; i < len(b); i++ {
		token[i] = TOKEN_SOURCE_STRING[int(b[i])%len(TOKEN_SOURCE_STRING)]
	}

	//将生成的token作为key,answer作为value存入Redis
	conn := l.svcCtx.RedisPool.Get()
	defer conn.Close()
	if _, err = conn.Do("SETEX", KEY_PREFIX_TOKEN+string(token[:]), EXPIRE_LENGHT, data.Answer); err != nil {
		return nil, err
	}

	return &types.QuestionReply{
		Question:     data.Question,
		RandomAnswer: string(randomAnswer),
		Token:        string(token[:]),
	}, nil
}

func RandomAnswer(answer string, disturbAnswer string) []byte {
	//多去少补
	rand.Seed(randomSeed)
	disturb := []byte(disturbAnswer)
	if len(answer+disturbAnswer)/3 > RANDOM_ANSWER_LENGTH {
		removeCount := len(answer+disturbAnswer)/3 - RANDOM_ANSWER_LENGTH
		remainCount := len(disturbAnswer)/3 - 1
		for i := 0; i < removeCount; i++ {
			index := rand.Intn(remainCount) * 3
			for j := 0; j < 3; j++ {
				disturb[index+j], disturb[remainCount*3+j] = disturb[remainCount*3+j], disturb[index+j]
			}
			remainCount--
		}
	} else {
		hanzi := RandomCommonHanzi(answer + disturbAnswer)
		disturb = append(disturb, []byte(hanzi)...)
	}
	//打乱顺序
	randomString := answer + string(disturb[:RANDOM_ANSWER_LENGTH*3-len(answer)])
	randomAnswer := []byte(randomString)
	for i := 0; i < len(randomAnswer)/3; i++ {
		r := rand.Intn(len(randomAnswer) / 3)
		for j := 0; j < 3; j++ {
			randomAnswer[i*3+j], randomAnswer[r*3+j] = randomAnswer[r*3+j], randomAnswer[i*3+j]
		}
	}
	randomSeed++

	return randomAnswer
}

// 获取常见汉字
func RandomCommonHanzi(hanzi string) string {
	if len(hanzi)/3 >= RANDOM_ANSWER_LENGTH {
		return hanzi
	}
	count := RANDOM_ANSWER_LENGTH - len(hanzi)/3
	rand.Seed(randomSeed)
	hanzis := make([][3]byte, count)
	checkHanzi := hanzi
	flag := false
	for i := 0; i < count; i++ {
		index := rand.Intn(len(COMMON_HANZI)/3-1) * 3
		copy(hanzis[i][:], COMMON_HANZI[index:index+3])
		for j := 0; j < len(checkHanzi)/3; j++ {
			if checkHanzi[j*3:j*3+3] == string(hanzis[i][:]) {
				i--
				flag = true
				break
			}
		}
		if !flag {
			checkHanzi += string(hanzis[i][:])
		}
		flag = false
	}
	randomSeed++
	s := ""
	for i := 0; i < count; i++ {
		s += string(hanzis[i][:])
	}
	return s
}

func UnicodetoHanzi() {
	sText := "\u7684\u4e00\u662f\u4e86\u6211"
	textQuoted := strconv.QuoteToASCII(sText)
	textUnquoted := textQuoted[1 : len(textQuoted)-1]
	sUnicodev := strings.Split(textUnquoted, "\\u")
	var context string
	for _, v := range sUnicodev {
		if len(v) < 1 {
			continue
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			panic(err)
		}
		context += fmt.Sprintf("%c", temp)
	}
}
