package main

import (
	"encoding/base64"
	"image/color"
	"math/big"
	"strconv"
	"strings"

	"hencryptor/icon"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	paillier "github.com/tao2581/go-go-gadget-paillier"
)

var privKey *paillier.PrivateKey

// Show loads a new text editor
func main() {
	app := app.NewWithID("cn.testbay.hencryptor")
	app.Settings().SetTheme(&myTheme{})
	window := app.NewWindow("Testbay HElib Encoder")
	window.SetIcon(icon.CalculatorBitmap)
	entry := widget.NewMultiLineEntry()
	output := widget.NewMultiLineEntry()
	output.Wrapping = fyne.TextWrapBreak
	entry.Wrapping = fyne.TextWrapBreak
	entry.SetPlaceHolder("输入数字")
	output.SetPlaceHolder("加密内容")

	// 加载key配置文件
	privKey = LoadKey(app)

	toolbar := fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		widget.NewButton("Encrypt", func() {
			// Encrypt
			lines := strings.Split(entry.Text, "\n")
			result := ""
			for _, value := range lines {
				if value == "" {
					continue
				}
				intValue, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					result += "[" + value + "] Error:输入的必须是数字"
				} else {
					valueInt64 := new(big.Int).SetInt64(intValue)
					ciper, _ := paillier.Encrypt(&privKey.PublicKey, valueInt64.Bytes())
					result += base64.StdEncoding.EncodeToString(ciper)
				}
				result += "\n"
			}
			output.SetText(result)

		}),
		widget.NewButton("Decrypt", func() {
			// Decrypt
			lines := strings.Split(output.Text, "\n")
			result := ""
			for _, value := range lines {
				if value == "" {
					continue
				}
				decodeText, _ := base64.StdEncoding.DecodeString(value)
				d, _ := paillier.Decrypt(privKey, decodeText)
				plainText := new(big.Int).SetBytes(d)
				result += plainText.String() + "\n"
			}
			entry.SetText(result)
		}),
		widget.NewButton("Clear", func() {
			entry.SetText("")
			output.SetText("")
		}),
		widget.NewButton("同态计算演示", func() {
			showExampleWindow(app, privKey)
		}),

		layout.NewSpacer(),
	)

	content := fyne.NewContainerWithLayout(
		layout.NewGridLayoutWithRows(2),
		fyne.NewContainerWithLayout(
			layout.NewBorderLayout(nil, toolbar, nil, nil),
			toolbar,
			widget.NewScrollContainer(entry),
		),
		widget.NewScrollContainer(output),
	)

	window.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Clear", func() {
				entry.SetText("")
			}),
		),

		fyne.NewMenu("秘钥管理",
			fyne.NewMenuItem("查看", func() {
				output.SetText("公钥:\n" + Pubkey2str(&privKey.PublicKey) + "\n\n私钥:\n" + Key2str(privKey) + "\n\n注意事项：\n请将公钥对应字符提交到轻流保存，私钥对应字符串务必复制到本地文件妥善保存备份，如私钥丢失线上所有加密数据都将无法找回")
			}),
			fyne.NewMenuItem("重置", func() {
				confirmHandler := func(confrm bool) {
					if confrm {
						privKey = NewKey(app)
						output.SetText("公钥:\n" + Pubkey2str(&privKey.PublicKey) + "\n\n私钥:\n" + Key2str(privKey) + "\n\n注意事项：请将公钥对应字符提交到轻流保存，私钥对应字符串务必复制到本地文件妥善保存备份，如私钥丢失线上所有加密数据都将无法找回")
					}
				}
				dialog.ShowConfirm("确认", "轻流中的加密数据将全部失效，确认要重置秘钥吗？", confirmHandler, window)
			}),
			fyne.NewMenuItem("恢复", func() {
				pk, err := RestoreKey(strings.Trim(entry.Text, " "), app)
				if err != nil {
					dialog.ShowError(err, window)
				} else {
					privKey = pk
					output.SetText("公钥:\n" + Pubkey2str(&privKey.PublicKey) + "\n\n私钥:\n" + Key2str(privKey) + "\n\n注意事项：恢复成功！\n请将公钥对应字符提交到轻流保存，私钥对应字符串务必复制到本地文件妥善保存备份，如私钥丢失线上所有加密数据都将无法找回")
				}

			}),
		),
	))
	window.SetContent(content)
	window.Resize(fyne.NewSize(800, 600))
	window.ShowAndRun()
}

func showExampleWindow(app fyne.App, pvKey *paillier.PrivateKey) {
	pubKey := pvKey.PublicKey
	var initCiper []byte

	secWindow := app.NewWindow("Calc Demo")
	label1 := canvas.NewText("1.加密公钥:", color.White)
	value1 := widget.NewMultiLineEntry()
	value1.Wrapping = fyne.TextWrapBreak
	value1.SetText(Pubkey2str(&pubKey))

	label2 := canvas.NewText("2. 原始数字", color.White)
	initEntry := widget.NewEntry()
	initEntry.SetText("15")
	labelEncryptLabel := canvas.NewText("3.加密后内容", color.White)
	labelEncryptText := widget.NewEntry()

	encryptBtn := widget.NewButton("加密", func() {
		intValue, err := strconv.ParseInt(initEntry.Text, 10, 64)
		if err != nil {
			dialog.ShowCustom("提示", "确定", canvas.NewText("输入的必须是数字", color.White), secWindow)
			return
		}
		m15 := new(big.Int).SetInt64(intValue)

		c15, _ := paillier.Encrypt(&pubKey, m15.Bytes())
		initCiper = c15
		labelEncryptText.SetText(base64.StdEncoding.EncodeToString(c15))
	})

	plusEntry := widget.NewEntry()
	plusEntry.SetText("10")
	plusResult := widget.NewLabel("")
	plusBtn := widget.NewButton("计算", func() {
		value, err := strconv.ParseInt(plusEntry.Text, 10, 64)
		if err != nil {
			dialog.ShowCustom("提示", "确定", canvas.NewText("输入的必须是数字", color.White), secWindow)
			return
		}

		plused := paillier.Add(&pubKey, initCiper, new(big.Int).SetInt64(value).Bytes())
		decrypted, _ := paillier.Decrypt(pvKey, plused)
		plusResult.SetText("\"" + base64.StdEncoding.EncodeToString(initCiper) + "\"  PLUS " +
			plusEntry.Text + " EQUAL: \"" + base64.StdEncoding.EncodeToString(plused) + "\" \n After decryption: " +
			new(big.Int).SetBytes(decrypted).String())
	})
	emptyLabel := canvas.NewText("", color.White)
	grid := fyne.NewContainerWithLayout(layout.NewFormLayout(),
		label1, value1, label2, initEntry,
		emptyLabel, encryptBtn,
		labelEncryptLabel, labelEncryptText,
		canvas.NewText("4. 同态加法", color.White), plusEntry,
		emptyLabel, plusBtn,
		canvas.NewText("结果", color.White), plusResult,
	)

	secWindow.SetContent(grid)
	secWindow.Resize(fyne.NewSize(600, 500))
	secWindow.Show()
}
