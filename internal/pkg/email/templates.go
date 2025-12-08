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
    <title>LexVeritas Verification</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f4f7f6; color: #333333;">
    <table role="presentation" style="width: 100%%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 100%%; max-width: 600px; background-color: #ffffff; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05); border: 1px solid #e1e4e8;">
                    
                    <!-- Header -->
                    <tr>
                        <td style="background-color: #111111; padding: 30px 40px; text-align: left;">
                            <h1 style="color: #ffffff; margin: 0; font-size: 24px; font-weight: 600; letter-spacing: 0.5px;">
                                LexVeritas
                            </h1>
                        </td>
                    </tr>

                    <!-- Body -->
                    <tr>
                        <td style="padding: 40px;">
                            <h2 style="margin: 0 0 24px; font-size: 20px; font-weight: 600; color: #111111;">邮箱验证</h2>
                            
                            <p style="margin: 0 0 24px; font-size: 15px; line-height: 1.6; color: #555555;">
                                您好，<br><br>
                                感谢您使用 LexVeritas。我们需要验证您的邮箱地址以继续操作。请在验证页面输入下方的验证码：
                            </p>

                            <div style="background-color: #f0f7ff; border-left: 4px solid #0056b3; border-radius: 4px; padding: 25px; margin: 30px 0; text-align: center;">
                                <span style="font-family: 'Courier New', Courier, monospace; font-size: 32px; font-weight: 700; color: #0056b3; letter-spacing: 4px; display: block;">%s</span>
                            </div>

                            <p style="margin: 0; font-size: 14px; line-height: 1.6; color: #666666;">
                                该验证码将在 <strong>%d 分钟</strong> 后验证失效。
                                <br>
                                如果这不是您的操作，请直接忽略此邮件，您的账户安全不会受到影响。
                            </p>
                        </td>
                    </tr>

                    <!-- Footer -->
                    <tr>
                        <td style="background-color: #f8f9fa; padding: 24px 40px; text-align: center; border-top: 1px solid #eeeeee;">
                            <p style="margin: 0; font-size: 12px; color: #999999; line-height: 1.5;">
                                © 2025 LexVeritas Inc. All rights reserved.
                                <br>
                                此为系统自动发送邮件，请勿回复。
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, code, expireMinutes)
}
