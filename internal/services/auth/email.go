package auth

var emailTemplate string = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
</head>
<body>
    <table cellpadding="0" cellspacing="0" style="margin: 0; padding: 0px; max-width: 600px; width: 100%; border-collapse: collapse; border: none;">
        <tbody>
            <tr class="header" style="max-width: 600px; width: 100%; position: relative; height: 168px; ">
                <td class="icon" style="width: 128px; max-width: 128px; position: relative; left: 0; height: 128px;">
                    <img src="https://gamification-api.s3.eu-central-1.amazonaws.com/favicon-128x128.png" alt="logo" border="0" width="128px" height="128px" style="display: block; width: 128px; padding-top: 20px; padding-left: 20px; padding-bottom: 20px;">
                </td>
                <td class="name" style="max-width: 300px; position: relative; right: 0; height: 128px; text-align: right;">
                    <p style="display: inline-block; max-width: 300px; width: 300px; font-family: Verdana, Geneva, sans-serif; font-size: 14px; text-align: right; color: #c0c0c0; position: relative; margin-right: 30px;">Система по учету и управлению достиженями сотрудников</p>
                </td>
            </tr>
            <tr>
                <td colspan="2">
                    <div style="display: flex; justify-content: center;"><hr style="width: 540px; height: 0px; border-color: #c0c0c0; border-spacing: 0; border-top: none; border-left: none; border-right: none;"></div>
                </td>
            </tr>
            <tr class="text-title" style="max-width: 600px; width: 100%;">
                <td colspan="2" style="font-family: Verdana, Geneva, sans-serif; font-size: 24px; max-width: 600px; width: 100%; padding-top: 20px; padding-left: 30px;">Здравствуйте!</td>
            </tr>
            <tr class="text-main" style="max-width: 600px; width: 100%;">
                <td colspan="2" style="font-family: Verdana, Geneva, sans-serif; font-size: 18px; max-width: 600px; width: 100%; padding-top: 20px; padding-left: 30px;">Для завершения авторизации в системе, введите на странице подтверждения код, представленный ниже:</td>
            </tr>
            <tr class="code">
                <td colspan="2" style="position: relative; height: 60px; padding-bottom: 30px;">
                    <div style=" width:540px;">
                        <div style="width: 540px; height: 60px; background: #1976d2; border-radius: 10px; margin-top: 30px; margin-left: 30px;">
                            <p style="font-family: Verdana, Geneva, sans-serif; font-size: 30px; font-weight: bold; margin: 0px; line-height: 60px; color: #ffffff; text-align: center;">{{.Code}}</p>
                        </div>
                    </div>
                </td>
            </tr>
            <tr>
                <td colspan="2">
                    <div style="display: flex; justify-content: center;"><hr style="width: 540px; height: 0px; border-color: #c0c0c0; border-spacing: 0; border-top: none; border-left: none; border-right: none;"></div>
                </td>
            </tr>
            <tr class="footer" style="max-width: 600px; width: 100%;">
                <td colspan="2">
                    <p style="font-family: Verdana, Geneva, sans-serif; font-size: 14px; color: #c0c0c0; padding-left: 30px;">Если Вы получили это письмо по ошибке, просто удалите его.</p>
                    <p style="font-family: Verdana, Geneva, sans-serif; font-size: 14px; color: #c0c0c0; padding-left: 30px; padding-top: 20px;">&#169; Achieveit</p>
                </td>
            </tr>
        </tbody>
    </table>
</body>
</html>`
