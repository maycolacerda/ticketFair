// services/email_templates.go
package services

import "fmt"

type emailTemplate struct {
	Subject string
	Body    string
}

func verificationEmailTemplate(code, expiryMinutes string) emailTemplate {
	return emailTemplate{
		Subject: "TicketFair — Código de verificação",
		Body: fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body        { font-family: Arial, sans-serif; background: #f4f4f4; margin: 0; padding: 0; }
    .container  { max-width: 600px; margin: 40px auto; background: #ffffff; border-radius: 8px; overflow: hidden; }
    .header     { background: #1a1a2e; padding: 32px; text-align: center; }
    .header h1  { color: #ffffff; margin: 0; font-size: 24px; }
    .body       { padding: 40px 32px; }
    .code-box   { background: #f0f4ff; border: 2px dashed #4361ee; border-radius: 8px;
                  text-align: center; padding: 24px; margin: 32px 0; }
    .code       { font-size: 48px; font-weight: bold; letter-spacing: 12px; color: #1a1a2e; }
    .expiry     { color: #888; font-size: 14px; margin-top: 8px; }
    .footer     { background: #f8f8f8; padding: 24px; text-align: center; color: #aaa; font-size: 12px; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>🎟️ TicketFair</h1>
    </div>
    <div class="body">
      <h2>Verificação de email</h2>
      <p>Use o código abaixo para verificar seu endereço de email.</p>
      <div class="code-box">
        <div class="code">%s</div>
        <div class="expiry">Expira em %s minutos</div>
      </div>
      <p>Se você não solicitou este código, ignore este email.</p>
    </div>
    <div class="footer">
      &copy; 2026 TicketFair. Todos os direitos reservados.
    </div>
  </div>
</body>
</html>
`, code, expiryMinutes),
	}
}

func welcomeEmailTemplate(username string) emailTemplate {
	return emailTemplate{
		Subject: "Bem-vindo ao TicketFair! 🎟️",
		Body: fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body        { font-family: Arial, sans-serif; background: #f4f4f4; margin: 0; padding: 0; }
    .container  { max-width: 600px; margin: 40px auto; background: #ffffff; border-radius: 8px; overflow: hidden; }
    .header     { background: #1a1a2e; padding: 32px; text-align: center; }
    .header h1  { color: #ffffff; margin: 0; font-size: 24px; }
    .body       { padding: 40px 32px; }
    .btn        { display: inline-block; background: #4361ee; color: #ffffff; padding: 14px 32px;
                  border-radius: 6px; text-decoration: none; font-weight: bold; margin-top: 24px; }
    .footer     { background: #f8f8f8; padding: 24px; text-align: center; color: #aaa; font-size: 12px; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>🎟️ TicketFair</h1>
    </div>
    <div class="body">
      <h2>Olá, %s! 👋</h2>
      <p>Seja bem-vindo ao TicketFair! Sua conta foi criada com sucesso.</p>
      <p>Agora você pode:</p>
      <ul>
        <li>Descobrir eventos próximos</li>
        <li>Comprar ingressos com segurança</li>
        <li>Gerenciar seus ingressos em um só lugar</li>
      </ul>
      <a href="http://ticketfair.localhost" class="btn">Explorar eventos</a>
    </div>
    <div class="footer">
      &copy; 2026 TicketFair. Todos os direitos reservados.
    </div>
  </div>
</body>
</html>
`, username),
	}
}

func purchaseConfirmationTemplate(username, eventName, ticketID string, amount float64) emailTemplate {
	return emailTemplate{
		Subject: fmt.Sprintf("Ingresso confirmado — %s 🎟️", eventName),
		Body: fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body        { font-family: Arial, sans-serif; background: #f4f4f4; margin: 0; padding: 0; }
    .container  { max-width: 600px; margin: 40px auto; background: #ffffff; border-radius: 8px; overflow: hidden; }
    .header     { background: #1a1a2e; padding: 32px; text-align: center; }
    .header h1  { color: #ffffff; margin: 0; font-size: 24px; }
    .body       { padding: 40px 32px; }
    .ticket     { background: #f0f4ff; border-left: 4px solid #4361ee; border-radius: 4px;
                  padding: 20px 24px; margin: 24px 0; }
    .ticket h3  { margin: 0 0 12px 0; color: #1a1a2e; }
    .ticket p   { margin: 4px 0; color: #555; font-size: 14px; }
    .footer     { background: #f8f8f8; padding: 24px; text-align: center; color: #aaa; font-size: 12px; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>🎟️ TicketFair</h1>
    </div>
    <div class="body">
      <h2>Compra confirmada! ✅</h2>
      <p>Olá %s, seu ingresso foi confirmado.</p>
      <div class="ticket">
        <h3>%s</h3>
        <p><strong>ID do ingresso:</strong> %s</p>
        <p><strong>Valor pago:</strong> R$ %.2f</p>
      </div>
      <p>Apresente o ID do ingresso na entrada do evento.</p>
    </div>
    <div class="footer">
      &copy; 2026 TicketFair. Todos os direitos reservados.
    </div>
  </div>
</body>
</html>
`, username, eventName, ticketID, amount),
	}
}
