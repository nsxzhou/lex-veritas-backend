package email

import "fmt"

// VerificationCodeTemplate 生成验证码邮件 HTML
func VerificationCodeTemplate(code string, expireMinutes int) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="margin: 0; padding: 0; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f4f4f4;">
    <table role="presentation" style="width: 100%%; max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; margin-top: 20px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);">
        <tr>
            <td style="padding: 40px 30px; text-align: center; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);">
                <h1 style="color: #ffffff; margin: 0; font-size: 28px;">LexVeritas</h1>
            </td>
        </tr>
        <tr>
            <td style="padding: 40px 30px;">
                <h2 style="color: #333333; margin: 0 0 20px 0; font-size: 22px;">邮箱验证</h2>
                <p style="color: #666666; font-size: 16px; line-height: 1.6; margin: 0 0 30px 0;">
                    您正在进行邮箱验证，请使用以下验证码完成验证：
                </p>
                <div style="background-color: #f8f9fa; border-radius: 8px; padding: 25px; text-align: center; margin-bottom: 30px;">
                    <span style="font-size: 36px; font-weight: bold; letter-spacing: 8px; color: #667eea;">%s</span>
                </div>
                <p style="color: #999999; font-size: 14px; line-height: 1.6; margin: 0;">
                    验证码有效期为 %d 分钟。如果您没有请求此验证码，请忽略此邮件。
                </p>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px 30px; background-color: #f8f9fa; text-align: center;">
                <p style="color: #999999; font-size: 12px; margin: 0;">
                    此邮件由系统自动发送，请勿回复。
                </p>
            </td>
        </tr>
    </table>
</body>
</html>
`, code, expireMinutes)
}
