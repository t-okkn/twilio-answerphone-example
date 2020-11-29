package main

import (
	"fmt"
	"time"

	"bitbucket.org/ckvist/twilio/twiml"
	"github.com/gin-gonic/gin"
	"github.com/t-okkn/go-enjaxel/notify"
)


// summary => 待ち受けるサーバのルーターを定義します
// return::*gin.Engine =>
// remark => httpHandlerを受け取る関数にそのまま渡せる
/////////////////////////////////////////
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// ルーティングの定義
	router.POST("/phonebot/v1/answerphone", responseAnswerPhone)

	 return router
}

// summary => 留守番電話応答時の動作を定義します
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func responseAnswerPhone(c *gin.Context) {
	resp := twiml.NewResponse()
	lang := "ja-jp"

	// POSTで"From"を取得
	from := c.PostForm("From")
	if from == "" {
		// なければBad Requestとみなす
		c.AbortWithStatus(400)
	}

	var msg string

	// 特殊番号（匿名）からの電話を拒否
	// c.f.) https://qiita.com/mobilebiz/items/e093f8bc3329114e7cd0
	if from == "+266696687" {
		msg = "非通知着信はお断りしております。"
		msg += "大変申し訳ありませんが発信番号を通知しておかけ直し下さい。"

		resp.Action(twiml.Say{ Text: msg, Language: lang })

	} else {
		msg = "ただ今電話に出ることができません。"
		msg += "ビープおんのあとにお名前とご用件をお願いいたします。"

		// 機械に喋らせたあとに録音時間を設ける
		resp.Action(
			twiml.Say{ Text: msg, Language: lang },
			twiml.Record{ Timeout: 120 },
		)
	}

	datetime := time.Now().Format("2006/01/02 15:04:05")
	// メッセージを送信する機能の関数ポインタを入れ込む
	sendmsg := sendLineMessage
	// sendmsg := sendSlackMessage
	// sendmsg := sendMail

	// 見やすいようにE.164形式から変換
	// （ただし、ハイフンは付きません）
	p := PhoneNumber(from)
	dfrom, err := p.To0ABJ()
	if err != nil {
		fmt.Println(err)
	}

	// 留守録が終わったあとにcallbackが来るので、それも拾えるように
	if c.PostForm("RecordingSid") != "" {
		vurl := c.PostForm("RecordingUrl") + ".mp3"
		subject := "【Twilio】留守電"
		body := fmt.Sprintf(
			"%s に %s からの着信において新規留守電が登録されました。\n\n%s",
			datetime,
			dfrom,
			vurl,
		)

		// メッセージを送信（本当は非同期にしたかった・・・修正するかも）
		if err := sendmsg(subject, body); err != nil {
			fmt.Println(err)
		}

	} else {
		subject := "【Twilio】着信あり"
		body := fmt.Sprintf(
			"%s に %s から着信がありました。",
			datetime,
			dfrom,
		)

		if err := sendmsg(subject, body); err != nil {
			fmt.Println(err)
		}
	}

	// c.XML() ←ginの機能でレスを返すとおかしくなるので、
	// gin.Context内のResponseWiterを使用してTwiMLを表示
	resp.Send(c.Writer)
}

// summary => LINE Notifyでメッセージを送ります
// param::subject => 件名
// param::body => 本文
// return::error => エラー
/////////////////////////////////////////
func sendLineMessage(subject, body string) error {
	token := "ABCDEFGHIJKLMNOPQabcdefghijklmnopq0123456789"
	line := notify.NewLineClientWithTag(token, subject)

	return line.SendMessage(body)
}

// summary => SlackのWebhookでメッセージを送ります
// param::subject => 件名
// param::body => 本文
// return::error => エラー
/////////////////////////////////////////
func sendSlackMessage(subject, body string) error {
	url := "https://hooks.slack.com/services/TXXXXXXXX/BXXXXXXXXXX/ABCDEFGHIJKLMN0123456789"
	s := notify.NewSlackMessageClient(url)

	str := subject + "\n\n" + body
	return s.SendSimpleMessage(str)
}

// summary => メールを送信します
// param::subject => 件名
// param::body => 本文
// return::error => エラー
/////////////////////////////////////////
func sendMail(subject, body string) error {
	s := notify.NewSmtpServer(
		"stmp.gmail.com",
		587,
		"username@gmail.com",
		"password",
	)

	return s.EasySendMail(
		"username@gmail.com",
		"hoge@example.com",
		subject,
		body,
	)
}

