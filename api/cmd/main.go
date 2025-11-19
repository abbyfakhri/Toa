package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abbyfakhri/toa-api/cmd/server"
	"github.com/abbyfakhri/toa-api/internal/services/email"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("no environment variable found")
	}

	emailClient, err := email.NewClient(email.EmailConfig{
		EmailFrom:     os.Getenv("EMAIL_FROM"),
		EmailPassword: os.Getenv("EMAIL_PASSWORD"),
		SMTPHost:      os.Getenv("SMTP_HOST"),
		SMTPPort:      os.Getenv("SMTP_PORT"),
		EmailAlias:    os.Getenv("EMAIL_ALIAS"),
	})

	if err != nil {
		panic(fmt.Sprintf("unable to connect to email client, err: %s", err.Error()))
	}

	// test(emailClient)

	defer emailClient.Quit()

	server := server.NewServer(server.ServerConfig{
		Port: os.Getenv("PORT"),
		Db: nil,
		EmailClient: emailClient,
	})

	_, err = server.Start()
	if err != nil {
		panic(fmt.Sprintf("unable to start http server, error: %s", err.Error()))
	}

	//	services.LoadServices(e, nil, emailClient)

}

func test(client email.EmailClient) error {
	err := client.SendMail(email.Email{
		To:      "abby.kuliah@gmail.com",
		Subject: "ini lagi nyoba sesuatu",
		Body:    "ini body nya yessss",
	})

	if err != nil {
		log.Print(err)
	}

	err = client.SendMail(email.Email{
		To:      "abby.kuliah@gmail.com",
		Subject: "ini lagi nyoba sesuatu",
		Template: `<!doctype html>
<html>
    <head>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>Rekafin</title>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
        <link
            href="https://fonts.googleapis.com/css2?family=Manrope:wght@200..800&display=swap"
            rel="stylesheet"
        />
        <style>
            * {
                margin: 0px;
                padding: 0px;
                font-family: "Manrope", sans-serif;
                font-optical-sizing: auto;
            }
            @media only screen and (max-width: 600px) {
                .content__wrapper {
                    max-width: 90% !important;
                    border-radius: 10px !important;
                }
                .padding__wrapper {
                    padding: 16px 20px !important;
                }
                .logo__size {
                    width: 64px !important;
                }
                .social__media__icon {
                    width: 16px !important;
                    height: 16px !important;
                }
                .main__content {
                    margin: 0px;
                    padding: 0px !important;
                }
                .content__title {
                    font-size: 18px !important;
                    padding: 0px !important;
                }
                .content__description {
                    font-size: 12px !important;
                }
                .list__wrapper {
                    margin-top: 12px !important;
                }
                .check__circle {
                    width: 16px !important;
                }
                .list__text {
                    font-size: 12px !important;
                }
                .cta__btn {
                    font-size: 14px !important;
                }
                .hr__style {
                    margin: 24px auto !important;
                }
                .copyright {
                    font-size: 12px !important;
                    font-weight: 400 !important;
                    margin-top: 16px !important;
                    padding-bottom: 32px !important;
                }
                .hidden {
                    display: none !important;
                }
                .d-block {
                    display: block !important;
                }
            }
            .text-description {
                font-size: 14px;
                line-height: 165%;
                letter-spacing: 0%;
                text-align: justify;
            }
            .text-benefit {
                font-size: 14px;
                line-height: 165%;
                letter-spacing: 0%;
                font-weight: 500;
                color: #2d2d2d;
            }
            .font-bold {
                font-weight: 700 !important;
            }
            span {
                font-weight: 700 !important;
            }
        </style>
    </head>
    <body style="margin: 0; padding: 0; background-color: #f4f9f8">
        <div
            class="main_container"
            style="
                width: 100%;
                background-color: #f4f9f8;
                font-family: &quot;Manrope&quot;, sans-serif;
                padding-top: 20px;
                padding-bottom: 20px;
                color: #2d2d2d;
            "
        >
            <img
                src="https://rekafin.id/email/rekafin.png"
                alt="rekafin.id"
                style="margin: 0px auto 20px; width: 100px; display: block"
            />
            <div
                class="content__wrapper"
                style="
                    max-width: 700px;
                    margin: auto;
                    background-color: white;
                    border-radius: 15px;
                    width: 100%;
                "
            >
                <div
                    class="padding__wrapper"
                    style="padding: 30px 20px; margin: auto"
                >
                    <div
                        style="padding: 0px 22px; max-width: 580px"
                        class="main__content"
                    >
                        <img
                            src="https://rekafin.id/email/rekafin-mail.png"
                            alt="rekasale"
                            style="width: 100%; margin: auto; display: block"
                        />
                        <p
                            style="
                                color: #3a519d;
                                font-weight: 800;
                                text-align: center;
                                margin-top: 24px;
                                font-size: 18px;
                                line-height: 150%;
                                letter-spacing: 0%;
                            "
                        >
                            Rekafin - Tools Catat
                            <br class="d-block" style="display: none" />
                            Keuangan Berbasis AI di WhatsApp
                        </p>
                        <p
                            style="
                                margin-top: 12px;
                                line-height: 165%;
                                padding: 0;
                                font-size: 14px;
                                text-align: justify;
                            "
                        >
                            Punya banyak pengeluaran tapi ribet kalau harus
                            catat manual? Sekarang saatnya beralih ke
                            <span style="font-weight: 700">
                                Rekafin, produk original dari PT Seakun Global
                                Teknologi</span
                            >
                            yang bikin pencatatan keuangan jadi super mudah,
                            langsung lewat
                            <span style="font-weight: 700"
                                >WhatsApp dan berbasis AI.</span
                            >
                        </p>
                        <p
                            style="
                                margin-top: 4px;
                                line-height: 165%;
                                padding: 0;
                                font-size: 14px;
                                text-align: justify;
                            "
                        >
                            Nikmati beragam benefit diantaranya:
                        </p>
                        <table
                            cellpadding="0"
                            cellspacing="0"
                            border="0"
                            style="border-collapse: collapse; margin-top: 12px"
                        >
                            <tr>
                                <td
                                    style="
                                        vertical-align: top;
                                        padding-right: 6px;
                                    "
                                >
                                    <img
                                        src="https://rekafin.id/email/check-circle.png"
                                        alt="checked"
                                        width="20"
                                        height="20"
                                        style="display: block; margin-top: 1px"
                                    />
                                </td>
                                <td class="text-benefit">
                                    Catat pengeluaran & pemasukan langsung dari
                                    WhatsApp, gak perlu buka aplikasi
                                </td>
                            </tr>
                            <tr>
                                <td colspan="2" height="6"></td>
                            </tr>
                            <tr style="margin-top: 4px">
                                <td
                                    style="
                                        vertical-align: top;
                                        padding-right: 6px;
                                    "
                                >
                                    <img
                                        src="https://rekafin.id/email/check-circle.png"
                                        alt="checked"
                                        width="20"
                                        height="20"
                                        style="display: block; margin-top: 1px"
                                    />
                                </td>
                                <td class="text-benefit">
                                    Proses pencatatan keuangan lebih mudah tanpa
                                    perintah baku karena berbasis AI
                                </td>
                            </tr>
                            <tr>
                                <td colspan="2" height="6"></td>
                            </tr>
                            <tr style="margin-top: 4px">
                                <td
                                    style="
                                        vertical-align: top;
                                        padding-right: 6px;
                                    "
                                >
                                    <img
                                        src="https://rekafin.id/email/check-circle.png"
                                        alt="checked"
                                        width="20"
                                        height="20"
                                        style="display: block; margin-top: 1px"
                                    />
                                </td>
                                <td class="text-benefit">
                                    Rekap keuangan dengan upload foto
                                    bill/kwitansi
                                </td>
                            </tr>
                            <tr>
                                <td colspan="2" height="6"></td>
                            </tr>
                            <tr style="margin-top: 4px">
                                <td
                                    style="
                                        vertical-align: top;
                                        padding-right: 6px;
                                    "
                                >
                                    <img
                                        src="https://rekafin.id/email/check-circle.png"
                                        alt="checked"
                                        width="20"
                                        height="20"
                                        style="display: block; margin-top: 1px"
                                    />
                                </td>
                                <td class="text-benefit">
                                    Pencatatan keuangan terekap secara realtime
                                </td>
                            </tr>
                            <tr>
                                <td colspan="2" height="6"></td>
                            </tr>
                            <tr style="margin-top: 4px">
                                <td
                                    style="
                                        vertical-align: top;
                                        padding-right: 6px;
                                    "
                                >
                                    <img
                                        src="https://rekafin.id/email/check-circle.png"
                                        alt="checked"
                                        width="20"
                                        height="20"
                                        style="display: block; margin-top: 1px"
                                    />
                                </td>
                                <td class="text-benefit">
                                    Cek mutasi harian, mingguan atau bulanan
                                    langsung dari WhatsApp
                                </td>
                            </tr>
                            <tr>
                                <td colspan="2" height="6"></td>
                            </tr>
                            <tr style="margin-top: 4px">
                                <td
                                    style="
                                        vertical-align: top;
                                        padding-right: 6px;
                                    "
                                >
                                    <img
                                        src="https://rekafin.id/email/check-circle.png"
                                        alt="checked"
                                        width="20"
                                        height="20"
                                        style="display: block; margin-top: 1px"
                                    />
                                </td>
                                <td class="text-benefit">
                                    Analisa pemasukan & pengeluaran langsung
                                    dari WhatsApp
                                </td>
                            </tr>
                            <tr>
                                <td colspan="2" height="6"></td>
                            </tr>
                            <tr style="margin-top: 4px">
                                <td
                                    style="
                                        vertical-align: top;
                                        padding-right: 6px;
                                    "
                                >
                                    <img
                                        src="https://rekafin.id/email/check-circle.png"
                                        alt="checked"
                                        width="20"
                                        height="20"
                                        style="display: block; margin-top: 1px"
                                    />
                                </td>
                                <td class="text-benefit">
                                    Dilengkapi dashboard dengan banyak fitur
                                    untuk analisa
                                </td>
                            </tr>
                        </table>
                        <p class="text-description" style="margin-top: 12px">
                            Cocok untuk mengatur dan mengelola keuangan
                            <span style="font-weight: 700">Personal</span> dan
                            <span style="font-weight: 700">Bisnis.</span>
                        </p>
                        <a
                            href="https://rekafin.id"
                            target="_blank"
                            style="text-decoration: none"
                        >
                            <div
                                class="cta-btn"
                                style="
                                    background-color: #3a519d;
                                    border-radius: 8px;
                                    padding: 16px 56px;
                                    border: none;
                                    color: #ffffff;
                                    font-size: 14px;
                                    font-weight: 700;
                                    cursor: pointer;
                                    text-decoration: none;
                                    max-width: 580px;
                                    margin: 0 auto;
                                    margin-top: 20px;
                                    text-align: center;
                                    box-sizing: border-box;
                                "
                            >
                                Langganan Sekarang
                            </div>
                        </a>
                        <hr
                            class="hr__style"
                            style="
                                border: 0;
                                border-top: 1px solid #eff4f3;
                                width: 100%;
                                max-width: 398px;
                                margin: 24px auto;
                            "
                        />
                        <p
                            style="
                                text-align: center;
                                margin: 0;
                                padding: 0;
                                line-height: 20px;
                                letter-spacing: 0px;
                                font-size: 12px;
                                vertical-align: middle;
                            "
                            class="instruction"
                        >
                            E-mail ini dibuat otomatis, mohon tidak membalas.
                        </p>
                    </div>
                </div>
            </div>
            <div style="text-align: center">
                <p
                    class="copyright"
                    style="
                        color: #3a519d;
                        margin-top: 20px;
                        font-size: 14px;
                        padding-bottom: 20px;
                        letter-spacing: 0px;
                        font-weight: 500;
                    "
                >
                    Copyright 2025 @ rekafin.id All Right Reserved
                </p>
            </div>
        </div>
    </body>
</html>
`,
	})

	return err
}
