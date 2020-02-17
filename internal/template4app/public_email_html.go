package template4app

var (
	PublicAlertNotification = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns="http://www.w3.org/1999/xhtml">
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<meta name="viewport" content="width=device-width" />

<style>body {
width: 100% !important; min-width: 100%; -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; margin: 0; padding: 0;
}
img {
outline: none; text-decoration: none; -ms-interpolation-mode: bicubic; width: auto; float: left; clear: both; display: block;
}
body {
color: #222222; font-family: "Helvetica", "Arial", sans-serif; font-weight: normal; padding: 0; margin: 0; text-align: left; line-height: 1.3;
}
body {
font-size: 14px; line-height: 19px;
}
a:hover {
color: #2795b6 !important;
}
a:active {
color: #2795b6 !important;
}
a:visited {
color: #2ba6cb !important;
}
body {
font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none;
}
a:hover {
color: #ff8f2b !important;
}
a:active {
color: #F2821E !important;
}
a:visited {
color: #E67612 !important;
}
.better-button:hover a {
color: #FFFFFF !important; background-color: #F2821E; border: 1px solid #F2821E;
}
.better-button:visited a {
color: #FFFFFF !important;
}
.better-button:active a {
color: #FFFFFF !important;
}
.better-button-alt:hover a {
color: #ff8f2b !important; background-color: #DDDDDD; border: 1px solid #F2821E;
}
.better-button-alt:visited a {
color: #ff8f2b !important;
}
.better-button-alt:active a {
color: #ff8f2b !important;
}
body {
height: 100% !important; width: 100% !important;
}
body .copy {
-ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;
}
.ExternalClass {
width: 100%;
}
.ExternalClass {
line-height: 100%;
}
img {
-ms-interpolation-mode: bicubic;
}
img {
border: 0 !important; outline: none !important; text-decoration: none !important;
}
a:hover {
text-decoration: underline;
}
@media only screen and (max-width: 600px) {
  table[class="body"] center {
    min-width: 0 !important;
  }
  table[class="body"] .container {
    width: 95% !important;
  }
  table[class="body"] .row {
    width: 100% !important; display: block !important;
  }
  table[class="body"] .wrapper {
    display: block !important; padding-right: 0 !important;
  }
  table[class="body"] .columns {
    table-layout: fixed !important; float: none !important; width: 100% !important; padding-right: 0px !important; padding-left: 0px !important; display: block !important;
  }
  table[class="body"] table.columns td {
    width: 100% !important;
  }
  table[class="body"] .columns td.six {
    width: 50% !important;
  }
  table[class="body"] .columns td.twelve {
    width: 100% !important;
  }
  table[class="body"] table.columns td.expander {
    width: 1px !important;
  }
  .logo {
    margin-left: 10px;
  }
}
@media (max-width: 600px) {
  table[class="email-container"] {
    width: 95% !important;
  }
  img[class="fluid"] {
    width: 100% !important; max-width: 100% !important; height: auto !important; margin: auto !important;
  }
  img[class="fluid-centered"] {
    width: 100% !important; max-width: 100% !important; height: auto !important; margin: auto !important;
  }
  img[class="fluid-centered"] {
    margin: auto !important;
  }
  td[class="comms-content"] {
    padding: 20px !important;
  }
  td[class="stack-column"] {
    display: block !important; width: 100% !important; direction: ltr !important;
  }
  td[class="stack-column-center"] {
    display: block !important; width: 100% !important; direction: ltr !important;
  }
  td[class="stack-column-center"] {
    text-align: center !important;
  }
  td[class="copy"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="copy -center"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="copy -bold"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="small-text"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="mini-centered-text"] {
    font-size: 14px !important; line-height: 24px !important; padding: 15px 30px !important;
  }
  td[class="copy -padd"] {
    padding: 0 40px !important;
  }
  span[class="sep"] {
    display: none !important;
  }
  td[class="mb-hide"] {
    display: none !important; height: 0 !important;
  }
  td[class="spacer mb-shorten"] {
    height: 25px !important;
  }
  .two-up td {
    width: 270px;
  }
}
</style></head>
<body leftmargin="0" topmargin="0" marginwidth="0" marginheight="0" class="main" style="height: 100% !important; width: 100% !important; min-width: 100%; -webkit-text-size-adjust: none; -ms-text-size-adjust: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; text-align: left; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; margin: 0 auto; padding: 0;" bgcolor="#2e2e2e">

	<table class="body" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; height: 100%; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" bgcolor="#2e2e2e">
		<tr style="vertical-align: top; padding: 0;" align="left">
			<td class="center" align="center" valign="top" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;">
        <center style="width: 100%; min-width: 580px;">
					<table class="row header" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; margin-top: 25px; margin-bottom: 25px; padding: 0px;">
						<tr style="vertical-align: top; padding: 0;" align="left">
						  <td class="center" align="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" valign="top">
						    <center style="width: 100%; min-width: 580px;">

						      <table class="container" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: inherit; width: 580px; margin: 0 auto; padding: 0;">
						        <tr style="vertical-align: top; padding: 0;" align="left">
						          <td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">

						            <table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
						              <tr style="vertical-align: top; padding: 0;" align="left">
						                <td class="twelve sub-columns center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; min-width: 0px; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 10px 10px 0px;" align="center" valign="top">
                              <img class="logo" src="http://grafana.org/assets/img/logo_new_transparent_200x48.png" style="width: 200px; display: inline; outline: none !important; text-decoration: none !important; -ms-interpolation-mode: bicubic; clear: both; border: 0;" align="none" />
                            </td>
                            <td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
                          </tr>
						            </table>

						          </td>
						        </tr>
						      </table>

						    </center>
						  </td>
						</tr>
					</table>

					<table class="container" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: inherit; width: 580px; margin: 0 auto; padding: 0;" width="600" bgcolor="#efefef">
						<tr style="vertical-align: top; padding: 0;" align="left">
							<td height="2" class="spacer mb-shorten" style="font-size: 0; line-height: 0; mso-table-lspace: 0pt; mso-table-rspace: 0pt; background-image: linear-gradient(to right, #ffed00 0%, #f26529 75%); height: 2px !important; word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0; border: 0;" valign="top" align="left"> </td>
						</tr>
						<tr style="vertical-align: top; padding: 0;" align="left">
							<td class="mini-centered-text" style="color: #343b41; mso-table-lspace: 0pt; mso-table-rspace: 0pt; word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 25px 35px; font: 400 16px/27px 'Helvetica Neue', Helvetica, Arial, sans-serif;" align="center" valign="top">
								\{\{Subject .Subject "\{\{.Title\}\}"\}\}

<table class="row" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; display: block; padding: 0px;">
  <tr style="vertical-align: top; padding: 0;" align="left">
    <td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">
      <table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
        <tr style="vertical-align: top; padding: 0;" align="left">
          <td class="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="center" valign="top">
            <h3 style="/*text-align: center*/; color: \{\{.SeverityColor\}\}; font-weight: bold; font-style: italic; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; line-height: 1.3; word-break: normal; font-size: 22px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 10px 0; padding: 0;" align="left">\{\{.Title\}\}</h3>
          </td>
        </tr>
      </table>
    </td>
  </tr>
</table>

<table class="row" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; display: block; padding: 0px;">
  <tr style="vertical-align: top; padding: 0;" align="left">
    <td class="last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0 0px 0 0;" align="left" valign="top">
      <table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
        <tr style="vertical-align: top; padding: 0;" align="left">
          <td class="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="center" valign="top">
            <p style="/*text-align: center*/; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0 0 10px; padding: 0;" align="left">\{\{.Message\}\}</p>
          </td>
        </tr>
      </table>
    </td>
  </tr>
</table>

\{\{if ne .Error "" \}\}
<table class="row" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; display: block; padding: 0px;">
  <tr style="vertical-align: top; padding: 0;" align="left">
    <td class="last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0 0px 0 0;" align="left" valign="top">
      <center style="width: 100%; min-width: 580px;">
      <table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
        <tr style="vertical-align: top; padding: 0;" align="left">
          <td class="twelve last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="left" valign="top">
            <h5 style="font-weight: bold; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; line-height: 1.3; word-break: normal; font-size: 18px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left">Error message</h5>
          </td>
        </tr>
        <tr style="vertical-align: top; padding: 0;" align="left">
          <td class="twelve last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="left" valign="top">
            <p style="color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0 0 10px; padding: 0;" align="left">\{\{.Error\}\}</p>
          </td>
        </tr>
      </table>
      </center>
    </td>
  </tr>
</table>
\{\{end\}\}

\{\{if ne .State "ok" \}\}
<table class="row" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; display: block; padding: 0px;">
  <tr style="vertical-align: top; padding: 0;" align="left">
    <td class="last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0 0px 0 0;" align="left" valign="top">
      <center style="width: 100%; min-width: 580px;">
      <table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
        <tr style="vertical-align: top; padding: 0;" align="left">
          <td class="six" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; width: 50%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="left" valign="top">
            <h5 style="font-weight: bold; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; line-height: 1.3; word-break: normal; font-size: 18px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left">Metric name</h5>
          </td>
          <td class="six last" style="width: 100px; word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="right" valign="top">
            <h5 style="font-weight: bold; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; line-height: 1.3; word-break: normal; font-size: 18px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="right">Value</h5>
          </td>
        </tr>
        \{\{range .EvalMatches\}\}
        <tr style="vertical-align: top; padding: 0;" align="left">
          <td class="six" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; width: 50%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="left" valign="top">
            <h5 class="data" style="color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 1.3; word-break: normal; font-size: 16px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left">\{\{.Metric\}\}</h5>
          </td>
          <td class="six last" style="width: 100px; word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="right" valign="top">
            <h5 class="data" style="color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 1.3; word-break: normal; font-size: 16px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="right">\{\{.Value\}\}</h5>
          </td>
        </tr>
        \{\{end\}\}
      </table>
      </center>
    </td>
  </tr>
</table>
\{\{end\}\}

<table class="row" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; display: block; padding: 0px;">
    <tr style="vertical-align: top; padding: 0;" align="left">
    <td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">
      <table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
        <tr style="vertical-align: top; padding: 0;" align="left">
          <td class="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="center" valign="top">
            \{\{if ne .ImageLink "" \}\}
              <img src="\{\{.ImageLink\}\}" alt="Alerting Panel" style="outline: none !important; text-decoration: none !important; -ms-interpolation-mode: bicubic; width: auto; clear: both; display: block; border: 0;" align="left" />
            \{\{end\}\}
            \{\{if ne .EmbeddedImage "" \}\}
              <img src="cid:\{\{.EmbeddedImage\}\}" alt="Alerting Panel" style="outline: none !important; text-decoration: none !important; -ms-interpolation-mode: bicubic; width: auto; clear: both; display: block; border: 0;" align="left" />
            \{\{end\}\}
          </td>
        </tr>
      </table>
    </td>
  </tr>
</table>


<table class="row" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; display: block; padding: 0px;">
  <tr style="vertical-align: top; padding: 0;" align="left">
    <td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">
      <table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
        <tr style="vertical-align: top; padding: 0;" align="left">
          <td class="center six" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; width: 50%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="center" valign="top">
            <table class="better-button" align="center" border="0" cellspacing="0" cellpadding="0" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; margin-top: 10px; margin-bottom: 20px; padding: 0;">
              <tr style="vertical-align: top; padding: 0;" align="left">
                <td align="center" class="better-button" bgcolor="#ff8f2b" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; -webkit-border-radius: 2px; -moz-border-radius: 2px; border-radius: 2px; margin: 0; padding: 0px;" valign="top">
                  <a href="\{\{.RuleUrl\}\}" target="_blank" style="color: #FFF; text-decoration: none; -webkit-border-radius: 2px; -moz-border-radius: 2px; border-radius: 2px; display: inline-block; padding: 12px 25px; border: 1px solid #ff8f2b;">View your Alert rule</a>
                </td>
              </tr>
            </table>
            </td>
            <td class="center six" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; width: 50%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="center" valign="top">
            <table class="better-button" align="center" border="0" cellspacing="0" cellpadding="0" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; margin-top: 10px; margin-bottom: 20px; padding: 0;">
              <tr style="vertical-align: top; padding: 0;" align="left">
                <td align="center" class="better-button-alt" bgcolor="#efefef" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; -webkit-border-radius: 2px; -moz-border-radius: 2px; border-radius: 2px; margin: 0; padding: 0px;" valign="top">
                  <a href="\{\{.AlertPageUrl\}\}" target="_blank" style="color: #ff8f2b; text-decoration: none; -webkit-border-radius: 2px; -moz-border-radius: 2px; border-radius: 2px; display: inline-block; background-color: #EFEFEF; padding: 12px 25px; border: 1px solid #ff8f2b;"> Go to the Alerts page</a>
                </td>
              </tr>
            </table>
          </td>
        </tr>
      </table>
    </td>
  </tr>
</table>




							</td>
						</tr>
					</table>

					<table class="footer center" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: center; color: #999999; margin-top: 20px; padding: 0;" bgcolor="#2e2e2e">
						<tr style="vertical-align: top; padding: 0;" align="left">
							<td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 20px 0px 0px;" align="left" valign="top">
								<table class="twelve columns center" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: center; width: 580px; margin: 0 auto; padding: 0;">
									<tr style="vertical-align: top; padding: 0;" align="left">
										<td class="twelve" align="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" valign="top">
											<center style="width: 100%; min-width: 580px;">
												<p style="font-size: 12px; color: #999999; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0 0 10px; padding: 0;" align="center">
													Sent by <a href="\{\{.AppUrl\}\}" style="color: #E67612; text-decoration: none;">Grafana v\{\{.BuildVersion\}\}</a>
													<br />© 2016 Grafana and raintank
												</p>
											</center>
										</td>
										<td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
									</tr>
								</table>
							</td>
						</tr>
					</table>
				</center>
			</td>
		</tr>
	</table>
</body>
</html>

`
	PublicAlertNotificationExample = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns="http://www.w3.org/1999/xhtml">
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<meta name="viewport" content="width=device-width" />
	
</head>
<body leftmargin="0" topmargin="0" marginwidth="0" marginheight="0" class="main" style="-ms-text-size-adjust: 100%; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; height: 100% !important; line-height: 19px; margin: 0 auto; min-width: 100%; padding: 0; text-align: left; width: 100% !important;" bgcolor="#2e2e2e"><style type="text/css">
body {
width: 100% !important; min-width: 100%; -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; margin: 0; padding: 0;
}
img {
outline: none; text-decoration: none; -ms-interpolation-mode: bicubic; width: auto; float: left; clear: both; display: block;
}
body {
color: #222222; font-family: "Helvetica", "Arial", sans-serif; font-weight: normal; padding: 0; margin: 0; text-align: left; line-height: 1.3;
}
body {
font-size: 14px; line-height: 19px;
}
a:hover {
color: #2795b6 !important;
}
a:active {
color: #2795b6 !important;
}
a:visited {
color: #2ba6cb !important;
}
body {
font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none;
}
a:hover {
color: #ff8f2b !important;
}
a:active {
color: #F2821E !important;
}
a:visited {
color: #E67612 !important;
}
.better-button:hover a {
color: #FFFFFF !important; background-color: #F2821E; border: 1px solid #F2821E;
}
.better-button:visited a {
color: #FFFFFF !important;
}
.better-button:active a {
color: #FFFFFF !important;
}
.better-button-alt:hover a {
color: #ff8f2b !important; background-color: #DDDDDD; border: 1px solid #F2821E;
}
.better-button-alt:visited a {
color: #ff8f2b !important;
}
.better-button-alt:active a {
color: #ff8f2b !important;
}
body {
height: 100% !important; width: 100% !important;
}
body .copy {
-ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;
}
.ExternalClass {
width: 100%;
}
.ExternalClass {
line-height: 100%;
}
img {
-ms-interpolation-mode: bicubic;
}
img {
border: 0 !important; outline: none !important; text-decoration: none !important;
}
a:hover {
text-decoration: underline;
}
@media only screen and (max-width: 600px) {
  table[class="body"] center {
    min-width: 0 !important;
  }
  table[class="body"] .container {
    width: 95% !important;
  }
  table[class="body"] .row {
    width: 100% !important; display: block !important;
  }
  table[class="body"] .wrapper {
    display: block !important; padding-right: 0 !important;
  }
  table[class="body"] .columns {
    table-layout: fixed !important; float: none !important; width: 100% !important; padding-right: 0px !important; padding-left: 0px !important; display: block !important;
  }
  table[class="body"] .column {
    table-layout: fixed !important; float: none !important; width: 100% !important; padding-right: 0px !important; padding-left: 0px !important; display: block !important;
  }
  table[class="body"] table.columns td {
    width: 100% !important;
  }
  table[class="body"] .columns td.six {
    width: 50% !important;
  }
  table[class="body"] .column td.six {
    width: 50% !important;
  }
  table[class="body"] .columns td.twelve {
    width: 100% !important;
  }
  table[class="body"] table.columns td.expander {
    width: 1px !important;
  }
  .logo {
    margin-left: 10px;
  }
}
@media (max-width: 600px) {
  table[class="email-container"] {
    width: 95% !important;
  }
  img[class="fluid"] {
    width: 100% !important; max-width: 100% !important; height: auto !important; margin: auto !important;
  }
  img[class="fluid-centered"] {
    width: 100% !important; max-width: 100% !important; height: auto !important; margin: auto !important;
  }
  img[class="fluid-centered"] {
    margin: auto !important;
  }
  td[class="comms-content"] {
    padding: 20px !important;
  }
  td[class="stack-column"] {
    display: block !important; width: 100% !important; direction: ltr !important;
  }
  td[class="stack-column-center"] {
    display: block !important; width: 100% !important; direction: ltr !important;
  }
  td[class="stack-column-center"] {
    text-align: center !important;
  }
  td[class="copy"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="copy -center"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="copy -bold"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="small-text"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="mini-centered-text"] {
    font-size: 14px !important; line-height: 24px !important; padding: 15px 30px !important;
  }
  td[class="copy -padd"] {
    padding: 0 40px !important;
  }
  span[class="sep"] {
    display: none !important;
  }
  td[class="mb-hide"] {
    display: none !important; height: 0 !important;
  }
  td[class="spacer mb-shorten"] {
    height: 25px !important;
  }
  .two-up td {
    width: 270px;
  }
}
</style>
	
	<table class="body" style="-webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; border-collapse: collapse; border-spacing: 0; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; height: 100%; line-height: 19px; margin: 0; padding: 0; text-align: left; vertical-align: top; width: 100%;">
		<tr style="padding: 0; vertical-align: top;" align="left">
			<td class="center" align="center" valign="top" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 0; word-break: break-word;">
        		<center style="min-width: 580px; width: 100%;">
					
					<table class="row header" style="border-collapse: collapse; border-spacing: 0; margin-bottom: 25px; margin-top: 25px; padding: 0px; position: relative; text-align: left; vertical-align: top; width: 100%;">
						<tr style="padding: 0; vertical-align: top;" align="left">
						  <td class="center" align="center" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 0; word-break: break-word;" valign="top">
						    <center style="min-width: 580px; width: 100%;">

						      <table class="container" style="border-collapse: collapse; border-spacing: 0; margin: 0 auto; padding: 0; text-align: inherit; vertical-align: top; width: 580px;">
						        <tr style="padding: 0; vertical-align: top;" align="left">
						          <td class="wrapper last" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 10px 0px 0px; position: relative; word-break: break-word;" align="left" valign="top">

						            <table class="twelve columns" style="border-collapse: collapse; border-spacing: 0; margin: 0 auto; padding: 0; text-align: left; vertical-align: top; width: 580px;">
						              <tr style="padding: 0; vertical-align: top;" align="left">
						                <td class="six sub-columns center" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; min-width: 0px; padding: 0px 10px 10px 0px; width: 50%; word-break: break-word;" align="center" valign="top">
											<img class="logo" src="http://grafana.org/assets/img/logo_new_transparent_200x48.png" style="-ms-interpolation-mode: bicubic; border: 0; clear: both; display: inline; outline: none !important; text-decoration: none !important; width: 200px;" align="none" />
						                </td>
										<td class="expander" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 0; visibility: hidden; width: 0px; word-break: break-word;" align="left" valign="top"></td>
						              </tr>
						            </table>

						          </td>
						        </tr>
						      </table>

						    </center>
						  </td>
						</tr>
					</table>

					<table class="container" style="border-collapse: collapse; border-spacing: 0; margin: 0 auto; padding: 0; text-align: inherit; vertical-align: top; width: 580px;" width="600" bgcolor="#efefef">
						<tr style="padding: 0; vertical-align: top;" align="left">
							<td height="2" class="spacer mb-shorten" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; background: linear-gradient(to right, #ffed00 0%, #f26529 75%); border: 0; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 0; font-weight: normal; height: 2px !important; hyphens: auto; line-height: 0; margin: 0; mso-table-lspace: 0pt; mso-table-rspace: 0pt; padding: 0; word-break: break-word;" valign="top" align="left"> </td>
						</tr>
						<tr style="padding: 0; vertical-align: top;" align="left">
							<td class="mini-centered-text" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #343b41; font: 400 16px/27px 'Helvetica Neue', Helvetica, Arial, sans-serif; hyphens: auto; margin: 0; mso-table-lspace: 0pt; mso-table-rspace: 0pt; padding: 25px 35px; word-break: break-word;" align="center" valign="top">
								<table class="row" style="border-collapse: collapse; border-spacing: 0; display: block; padding: 0px; position: relative; text-align: left; vertical-align: top; width: 100%;">
	<tr style="padding: 0; vertical-align: top;" align="left">
		<td class="wrapper last" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 10px 0px 0px; position: relative; word-break: break-word;" align="left" valign="top">
			<table class="twelve columns" style="border-collapse: collapse; border-spacing: 0; margin: 0 auto; padding: 0; text-align: left; vertical-align: top; width: 580px;">
				<tr style="padding: 0; vertical-align: top;" align="left">
					<td class="center" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 0px 0px 10px; word-break: break-word;" align="center" valign="top">
						<h3 style="-webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; /*text-align: center*/; color: #D63232; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 22px; font-style: italic; font-weight: bold; line-height: 1.3; margin: 20px 0 0; padding: 0; word-break: normal;" align="left"><img src="http://play.grafana.org/img/critical.svg" style="-ms-interpolation-mode: bicubic; border: 0; clear: both; display: block; margin-right: 10px; outline: none !important; text-decoration: none !important; width: 40px;" align="left" /> [CRITICAL] Imaginary timeserie alert</h3>
					</td>
				</tr>
			</table>
		</td>
	</tr>
</table>

<table class="row" style="border-collapse: collapse; border-spacing: 0; display: block; padding: 0px; position: relative; text-align: left; vertical-align: top; width: 100%;">
	<tr style="padding: 0; vertical-align: top;" align="left">
        <td class="wrapper last" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 10px 0px 0px; position: relative; word-break: break-word;" align="left" valign="top">
			<table class="twelve columns" style="border-collapse: collapse; border-spacing: 0; margin: 0 auto; padding: 0; text-align: left; vertical-align: top; width: 580px;">
				<tr style="padding: 0; vertical-align: top;" align="left">
                    <td class="center" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 0px 0px 10px; word-break: break-word;" align="center" valign="top">
                        <p style="-webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; /*text-align: center*/; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; line-height: 19px; margin: 0 0 10px; padding: 0;" align="left">Alert message that will support markdown in some distant future.</p>
					</td>
				</tr>
			</table>
		</td>
	</tr>
</table>


<table class="row" style="border-collapse: collapse; border-spacing: 0; display: block; padding: 0px; position: relative; text-align: left; vertical-align: top; width: 100%;">
	<tr style="padding: 0; vertical-align: top;" align="left">
        <td class="column last" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 10px 0px 0px; position: relative; word-break: break-word;" align="left" valign="top">
			<center style="min-width: 580px; width: 100%;">
			
			<table class="twelve columns" style="border-collapse: collapse; border-spacing: 0; margin: 0 auto; padding: 0; text-align: left; vertical-align: top; width: 580px;">
				<tr style="padding: 0; vertical-align: top;" align="left">
					<td class="six" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 18px; font-weight: bold; hyphens: auto; line-height: 1.3; margin: 0; padding: 0; width: 50%; word-break: normal;" align="left" valign="top">
						Metric name
					</td>
					<td class="six last" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 18px; font-weight: bold; hyphens: auto; line-height: 1.3; margin: 0; padding: 0; width: 50%; word-break: normal;" align="left" valign="top">
						Value
					</td>
				</tr>
				
				<tr style="padding: 0; vertical-align: top;" align="left">
					<td class="six" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 16px; font-weight: normal; hyphens: auto; line-height: 1.5; margin: 0; padding: 0; width: 50%; word-break: normal;" align="left" valign="top">
						desktop
					</td>
					<td class="six last" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 16px; font-weight: normal; hyphens: auto; line-height: 1.5; margin: 0; padding: 0; width: 50%; word-break: normal;" align="left" valign="top">
						40
					</td>
				</tr>
				
				<tr style="padding: 0; vertical-align: top;" align="left">
					<td class="six" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 16px; font-weight: normal; hyphens: auto; line-height: 1.5; margin: 0; padding: 0; width: 50%; word-break: normal;" align="left" valign="top">
						mobile
					</td>
					<td class="six last" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 16px; font-weight: normal; hyphens: auto; line-height: 1.5; margin: 0; padding: 0; width: 50%; word-break: normal;" align="left" valign="top">
						20
					</td>
				</tr>
				
			</table>
			</center>
		</td>
	</tr>
</table>


<table class="row" style="border-collapse: collapse; border-spacing: 0; display: block; padding: 0px; position: relative; text-align: left; vertical-align: top; width: 100%;">
    <tr style="padding: 0; vertical-align: top;" align="left">
		<td class="wrapper last" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 10px 0px 0px; position: relative; word-break: break-word;" align="left" valign="top">
			<table class="twelve columns" style="border-collapse: collapse; border-spacing: 0; margin: 0 auto; padding: 0; text-align: left; vertical-align: top; width: 580px;">
				<tr style="padding: 0; vertical-align: top;" align="left">
					<td class="center" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 0px 0px 10px; word-break: break-word;" align="center" valign="top">
	                    <img src="http://play.grafana.org/render/dashboard-solo/db/graphite-carbon-metrics?panelId=2&amp;from=1472149131057&amp;to=1472159931058&amp;width=1000&amp;height=500" style="-ms-interpolation-mode: bicubic; border: 0; clear: both; display: block; max-width: 100%; outline: none !important; text-decoration: none !important; width: auto;" align="left" />
					</td>
				</tr>
			</table>
		</td>
	</tr>
</table>


<table class="row" style="border-collapse: collapse; border-spacing: 0; display: block; padding: 0px; position: relative; text-align: left; vertical-align: top; width: 100%;">
	<tr style="padding: 0; vertical-align: top;" align="left">
		<td class="wrapper last" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 10px 0px 0px; position: relative; word-break: break-word;" align="left" valign="top">
			<table class="twelve columns" style="border-collapse: collapse; border-spacing: 0; margin: 0 auto; padding: 0; text-align: left; vertical-align: top; width: 580px;">
				<tr style="padding: 0; vertical-align: top;" align="left">
					<td class="center six" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 0px 0px 10px; width: 50%; word-break: break-word;" align="center" valign="top">
						<table class="better-button" align="center" border="0" cellspacing="0" cellpadding="0" style="border-collapse: collapse; border-spacing: 0; margin-bottom: 20px; margin-top: 10px; padding: 0; text-align: left; vertical-align: top;">
							<tr style="padding: 0; vertical-align: top;" align="left">
								<td align="center" class="better-button" bgcolor="#ff8f2b" style="-moz-border-radius: 2px; -moz-hyphens: auto; -webkit-border-radius: 2px; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; border-radius: 2px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 0px; word-break: break-word;" valign="top">
									<a href="\{\{.RuleUrl\}\}" target="_blank" style="-moz-border-radius: 2px; -webkit-border-radius: 2px; border: 1px solid #ff8f2b; border-radius: 2px; color: #FFF; display: inline-block; padding: 12px 25px; text-decoration: none;">View your Alert rule</a>
								</td>
							</tr>
						</table>
						</td>
						<td class="center six" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 0px 0px 10px; width: 50%; word-break: break-word;" align="center" valign="top">
						<table class="better-button" align="center" border="0" cellspacing="0" cellpadding="0" style="border-collapse: collapse; border-spacing: 0; margin-bottom: 20px; margin-top: 10px; padding: 0; text-align: left; vertical-align: top;">
							<tr style="padding: 0; vertical-align: top;" align="left">
								<td align="center" class="better-button-alt" bgcolor="#efefef" style="-moz-border-radius: 2px; -moz-hyphens: auto; -webkit-border-radius: 2px; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; border-radius: 2px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 0px; word-break: break-word;" valign="top">
									<a href="\{\{.AlertPageUrl\}\}" target="_blank" style="-moz-border-radius: 2px; -webkit-border-radius: 2px; background: #EFEFEF; border: 1px solid #ff8f2b; border-radius: 2px; color: #ff8f2b; display: inline-block; padding: 12px 25px; text-decoration: none;"> Go to the Alerts page</a>
								</td>
							</tr>
						</table>
					</td>
				</tr>
			</table>
		</td>
	</tr>
</table>


							
								
							</td>
						</tr>
					</table>
					
					<table class="footer center" style="border-collapse: collapse; border-spacing: 0; color: #999999; margin-top: 20px; padding: 0; text-align: center; vertical-align: top;" bgcolor="#2e2e2e">
						<tr style="padding: 0; vertical-align: top;" align="left">
							<td class="wrapper last" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 10px 20px 0px 0px; position: relative; word-break: break-word;" align="left" valign="top">
								<table class="twelve columns center" style="border-collapse: collapse; border-spacing: 0; margin: 0 auto; padding: 0; text-align: center; vertical-align: top; width: 580px;">
									<tr style="padding: 0; vertical-align: top;" align="left">
										<td class="twelve" align="center" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 0px 0px 10px; width: 100%; word-break: break-word;" valign="top">
											<center style="min-width: 580px; width: 100%;">
												<p style="-webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; color: #999999; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 12px; font-weight: normal; line-height: 19px; margin: 0 0 10px; padding: 0;" align="center">
													Sent by <a href="\{\{.AppUrl\}\}" style="color: #E67612; text-decoration: none;">Grafana v\{\{.BuildVersion\}\}</a>
													<br />© 2016 Grafana and raintank
												</p>
											</center>
										</td>
										<td class="expander" style="-moz-hyphens: auto; -webkit-font-smoothing: antialiased; -webkit-hyphens: auto; -webkit-text-size-adjust: none; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: normal; hyphens: auto; line-height: 19px; margin: 0; padding: 0; visibility: hidden; width: 0px; word-break: break-word;" align="left" valign="top"></td>
									</tr>
								</table>
							</td>
						</tr>
					</table>
				</center>
			</td>
		</tr>
	</table>
</body>
</html>

`
	PublicInvitedToOrg = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<meta name="viewport" content="width=device-width" />
	
<style>body {
width: 100% !important; min-width: 100%; -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; margin: 0; padding: 0;
}
img {
outline: none; text-decoration: none; -ms-interpolation-mode: bicubic; width: auto; float: left; clear: both; display: block;
}
body {
color: #222222; font-family: "Helvetica", "Arial", sans-serif; font-weight: normal; padding: 0; margin: 0; text-align: left; line-height: 1.3;
}
body {
font-size: 14px; line-height: 19px;
}
a:hover {
color: #2795b6 !important;
}
a:active {
color: #2795b6 !important;
}
a:visited {
color: #2ba6cb !important;
}
body {
font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none;
}
a:hover {
color: #ff8f2b !important;
}
a:active {
color: #F2821E !important;
}
a:visited {
color: #E67612 !important;
}
.better-button:hover a {
color: #FFFFFF !important; background-color: #F2821E; border: 1px solid #F2821E;
}
.better-button:visited a {
color: #FFFFFF !important;
}
.better-button:active a {
color: #FFFFFF !important;
}
.better-button-alt:hover a {
color: #ff8f2b !important; background-color: #DDDDDD; border: 1px solid #F2821E;
}
.better-button-alt:visited a {
color: #ff8f2b !important;
}
.better-button-alt:active a {
color: #ff8f2b !important;
}
body {
height: 100% !important; width: 100% !important;
}
body .copy {
-ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;
}
.ExternalClass {
width: 100%;
}
.ExternalClass {
line-height: 100%;
}
img {
-ms-interpolation-mode: bicubic;
}
img {
border: 0 !important; outline: none !important; text-decoration: none !important;
}
a:hover {
text-decoration: underline;
}
@media only screen and (max-width: 600px) {
  table[class="body"] center {
    min-width: 0 !important;
  }
  table[class="body"] .container {
    width: 95% !important;
  }
  table[class="body"] .row {
    width: 100% !important; display: block !important;
  }
  table[class="body"] .wrapper {
    display: block !important; padding-right: 0 !important;
  }
  table[class="body"] .columns {
    table-layout: fixed !important; float: none !important; width: 100% !important; padding-right: 0px !important; padding-left: 0px !important; display: block !important;
  }
  table[class="body"] table.columns td {
    width: 100% !important;
  }
  table[class="body"] .columns td.six {
    width: 50% !important;
  }
  table[class="body"] .columns td.twelve {
    width: 100% !important;
  }
  table[class="body"] table.columns td.expander {
    width: 1px !important;
  }
  .logo {
    margin-left: 10px;
  }
}
@media (max-width: 600px) {
  table[class="email-container"] {
    width: 95% !important;
  }
  img[class="fluid"] {
    width: 100% !important; max-width: 100% !important; height: auto !important; margin: auto !important;
  }
  img[class="fluid-centered"] {
    width: 100% !important; max-width: 100% !important; height: auto !important; margin: auto !important;
  }
  img[class="fluid-centered"] {
    margin: auto !important;
  }
  td[class="comms-content"] {
    padding: 20px !important;
  }
  td[class="stack-column"] {
    display: block !important; width: 100% !important; direction: ltr !important;
  }
  td[class="stack-column-center"] {
    display: block !important; width: 100% !important; direction: ltr !important;
  }
  td[class="stack-column-center"] {
    text-align: center !important;
  }
  td[class="copy"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="copy -center"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="copy -bold"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="small-text"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="mini-centered-text"] {
    font-size: 14px !important; line-height: 24px !important; padding: 15px 30px !important;
  }
  td[class="copy -padd"] {
    padding: 0 40px !important;
  }
  span[class="sep"] {
    display: none !important;
  }
  td[class="mb-hide"] {
    display: none !important; height: 0 !important;
  }
  td[class="spacer mb-shorten"] {
    height: 25px !important;
  }
  .two-up td {
    width: 270px;
  }
}
</style></head>
<body leftmargin="0" topmargin="0" marginwidth="0" marginheight="0" class="main" style="height: 100% !important; width: 100% !important; min-width: 100%; -webkit-text-size-adjust: none; -ms-text-size-adjust: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; text-align: left; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; margin: 0 auto; padding: 0;" bgcolor="#2e2e2e">

	<table class="body" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; height: 100%; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" bgcolor="#2e2e2e">
		<tr style="vertical-align: top; padding: 0;" align="left">
			<td class="center" align="center" valign="top" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;">
        <center style="width: 100%; min-width: 580px;">
					<table class="row header" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; margin-top: 25px; margin-bottom: 25px; padding: 0px;">
						<tr style="vertical-align: top; padding: 0;" align="left">
						  <td class="center" align="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" valign="top">
						    <center style="width: 100%; min-width: 580px;">

						      <table class="container" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: inherit; width: 580px; margin: 0 auto; padding: 0;">
						        <tr style="vertical-align: top; padding: 0;" align="left">
						          <td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">

						            <table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
						              <tr style="vertical-align: top; padding: 0;" align="left">
						                <td class="twelve sub-columns center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; min-width: 0px; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 10px 10px 0px;" align="center" valign="top">
                              <img class="logo" src="http://grafana.org/assets/img/logo_new_transparent_200x48.png" style="width: 200px; display: inline; outline: none !important; text-decoration: none !important; -ms-interpolation-mode: bicubic; clear: both; border: 0;" align="none" />
                            </td>
                            <td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
                          </tr>
						            </table>

						          </td>
						        </tr>
						      </table>

						    </center>
						  </td>
						</tr>
					</table>

					<table class="container" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: inherit; width: 580px; margin: 0 auto; padding: 0;" width="600" bgcolor="#efefef">
						<tr style="vertical-align: top; padding: 0;" align="left">
							<td height="2" class="spacer mb-shorten" style="font-size: 0; line-height: 0; mso-table-lspace: 0pt; mso-table-rspace: 0pt; background-image: linear-gradient(to right, #ffed00 0%, #f26529 75%); height: 2px !important; word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0; border: 0;" valign="top" align="left"> </td>
						</tr>
						<tr style="vertical-align: top; padding: 0;" align="left">
							<td class="mini-centered-text" style="color: #343b41; mso-table-lspace: 0pt; mso-table-rspace: 0pt; word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 25px 35px; font: 400 16px/27px 'Helvetica Neue', Helvetica, Arial, sans-serif;" align="center" valign="top">
								

\{\{Subject .Subject "\{\{.InvitedBy\}\} has added you to the \{\{.OrgName\}\} organization"\}\}

<table class="row" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; display: block; padding: 0px;">
	<tr style="vertical-align: top; padding: 0;" align="left">
		<td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">

			<table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
				<tr style="vertical-align: top; padding: 0;" align="left">
					<td style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="left" valign="top">
						<h4 class="center" style="color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 1.3; word-break: normal; font-size: 20px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="center">You have been added to \{\{.OrgName\}\}</h4>
					</td>
					<td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
				</tr>
			</table>

		</td>
	</tr>
</table>

<table class="row" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; display: block; padding: 0px;">
	<tr style="vertical-align: top; padding: 0;" align="left">
		<td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">
			<table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
				<tr style="vertical-align: top; padding: 0;" align="left">
					<td class="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="center" valign="top">
						<p style="color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0 0 10px; padding: 0;" align="left"><b>\{\{.InvitedBy\}\}</b> has added you to the <b>\{\{.OrgName\}\}</b> organization in Grafana.
						</p><p style="color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0 0 10px; padding: 0;" align="left">Once logged in, \{\{.OrgName\}\} will be available in the left side menu, in the dropdown below your username.</p>
					</td>
					<td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
				</tr>
				<tr style="vertical-align: top; padding: 0;" align="left">
					<td class="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="center" valign="top">
						<table class="better-button" align="center" border="0" cellspacing="0" cellpadding="0" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; margin-top: 10px; margin-bottom: 20px; padding: 0;">
							<tr style="vertical-align: top; padding: 0;" align="left">
								<td align="center" class="better-button" bgcolor="#ff8f2b" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; -webkit-border-radius: 2px; -moz-border-radius: 2px; border-radius: 2px; margin: 0; padding: 0px;" valign="top"><a href="\{\{.AppUrl\}\}" target="_blank" style="color: #FFF; text-decoration: none; -webkit-border-radius: 2px; -moz-border-radius: 2px; border-radius: 2px; display: inline-block; padding: 12px 25px; border: 1px solid #ff8f2b;">Log in now</a></td>
							</tr>
						</table>
					</td>
				</tr>
			</table>
		</td>
	</tr>
</table>



								
							</td>
						</tr>
					</table>
					
					<table class="footer center" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: center; color: #999999; margin-top: 20px; padding: 0;" bgcolor="#2e2e2e">
						<tr style="vertical-align: top; padding: 0;" align="left">
							<td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 20px 0px 0px;" align="left" valign="top">
								<table class="twelve columns center" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: center; width: 580px; margin: 0 auto; padding: 0;">
									<tr style="vertical-align: top; padding: 0;" align="left">
										<td class="twelve" align="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" valign="top">
											<center style="width: 100%; min-width: 580px;">
												<p style="font-size: 12px; color: #999999; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0 0 10px; padding: 0;" align="center">
													Sent by <a href="\{\{.AppUrl\}\}" style="color: #E67612; text-decoration: none;">Grafana v\{\{.BuildVersion\}\}</a>
													<br />© 2016 Grafana and raintank
												</p>
											</center>
										</td>
										<td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
									</tr>
								</table>
							</td>
						</tr>
					</table>
				</center>
			</td>
		</tr>
	</table>
</body>
</html>

`
	PublicNewUserInvite = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<meta name="viewport" content="width=device-width" />
	
<style>body {
width: 100% !important; min-width: 100%; -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; margin: 0; padding: 0;
}
img {
outline: none; text-decoration: none; -ms-interpolation-mode: bicubic; width: auto; float: left; clear: both; display: block;
}
body {
color: #222222; font-family: "Helvetica", "Arial", sans-serif; font-weight: normal; padding: 0; margin: 0; text-align: left; line-height: 1.3;
}
body {
font-size: 14px; line-height: 19px;
}
a:hover {
color: #2795b6 !important;
}
a:active {
color: #2795b6 !important;
}
a:visited {
color: #2ba6cb !important;
}
body {
font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none;
}
a:hover {
color: #ff8f2b !important;
}
a:active {
color: #F2821E !important;
}
a:visited {
color: #E67612 !important;
}
.better-button:hover a {
color: #FFFFFF !important; background-color: #F2821E; border: 1px solid #F2821E;
}
.better-button:visited a {
color: #FFFFFF !important;
}
.better-button:active a {
color: #FFFFFF !important;
}
.better-button-alt:hover a {
color: #ff8f2b !important; background-color: #DDDDDD; border: 1px solid #F2821E;
}
.better-button-alt:visited a {
color: #ff8f2b !important;
}
.better-button-alt:active a {
color: #ff8f2b !important;
}
body {
height: 100% !important; width: 100% !important;
}
body .copy {
-ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;
}
.ExternalClass {
width: 100%;
}
.ExternalClass {
line-height: 100%;
}
img {
-ms-interpolation-mode: bicubic;
}
img {
border: 0 !important; outline: none !important; text-decoration: none !important;
}
a:hover {
text-decoration: underline;
}
@media only screen and (max-width: 600px) {
  table[class="body"] center {
    min-width: 0 !important;
  }
  table[class="body"] .container {
    width: 95% !important;
  }
  table[class="body"] .row {
    width: 100% !important; display: block !important;
  }
  table[class="body"] .wrapper {
    display: block !important; padding-right: 0 !important;
  }
  table[class="body"] .columns {
    table-layout: fixed !important; float: none !important; width: 100% !important; padding-right: 0px !important; padding-left: 0px !important; display: block !important;
  }
  table[class="body"] table.columns td {
    width: 100% !important;
  }
  table[class="body"] .columns td.six {
    width: 50% !important;
  }
  table[class="body"] .columns td.twelve {
    width: 100% !important;
  }
  table[class="body"] table.columns td.expander {
    width: 1px !important;
  }
  .logo {
    margin-left: 10px;
  }
}
@media (max-width: 600px) {
  table[class="email-container"] {
    width: 95% !important;
  }
  img[class="fluid"] {
    width: 100% !important; max-width: 100% !important; height: auto !important; margin: auto !important;
  }
  img[class="fluid-centered"] {
    width: 100% !important; max-width: 100% !important; height: auto !important; margin: auto !important;
  }
  img[class="fluid-centered"] {
    margin: auto !important;
  }
  td[class="comms-content"] {
    padding: 20px !important;
  }
  td[class="stack-column"] {
    display: block !important; width: 100% !important; direction: ltr !important;
  }
  td[class="stack-column-center"] {
    display: block !important; width: 100% !important; direction: ltr !important;
  }
  td[class="stack-column-center"] {
    text-align: center !important;
  }
  td[class="copy"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="copy -center"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="copy -bold"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="small-text"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="mini-centered-text"] {
    font-size: 14px !important; line-height: 24px !important; padding: 15px 30px !important;
  }
  td[class="copy -padd"] {
    padding: 0 40px !important;
  }
  span[class="sep"] {
    display: none !important;
  }
  td[class="mb-hide"] {
    display: none !important; height: 0 !important;
  }
  td[class="spacer mb-shorten"] {
    height: 25px !important;
  }
  .two-up td {
    width: 270px;
  }
}
</style></head>
<body leftmargin="0" topmargin="0" marginwidth="0" marginheight="0" class="main" style="height: 100% !important; width: 100% !important; min-width: 100%; -webkit-text-size-adjust: none; -ms-text-size-adjust: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; text-align: left; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; margin: 0 auto; padding: 0;" bgcolor="#2e2e2e">

	<table class="body" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; height: 100%; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" bgcolor="#2e2e2e">
		<tr style="vertical-align: top; padding: 0;" align="left">
			<td class="center" align="center" valign="top" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;">
        <center style="width: 100%; min-width: 580px;">
					<table class="row header" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; margin-top: 25px; margin-bottom: 25px; padding: 0px;">
						<tr style="vertical-align: top; padding: 0;" align="left">
						  <td class="center" align="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" valign="top">
						    <center style="width: 100%; min-width: 580px;">

						      <table class="container" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: inherit; width: 580px; margin: 0 auto; padding: 0;">
						        <tr style="vertical-align: top; padding: 0;" align="left">
						          <td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">

						            <table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
						              <tr style="vertical-align: top; padding: 0;" align="left">
						                <td class="twelve sub-columns center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; min-width: 0px; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 10px 10px 0px;" align="center" valign="top">
                              <img class="logo" src="http://grafana.org/assets/img/logo_new_transparent_200x48.png" style="width: 200px; display: inline; outline: none !important; text-decoration: none !important; -ms-interpolation-mode: bicubic; clear: both; border: 0;" align="none" />
                            </td>
                            <td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
                          </tr>
						            </table>

						          </td>
						        </tr>
						      </table>

						    </center>
						  </td>
						</tr>
					</table>

					<table class="container" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: inherit; width: 580px; margin: 0 auto; padding: 0;" width="600" bgcolor="#efefef">
						<tr style="vertical-align: top; padding: 0;" align="left">
							<td height="2" class="spacer mb-shorten" style="font-size: 0; line-height: 0; mso-table-lspace: 0pt; mso-table-rspace: 0pt; background-image: linear-gradient(to right, #ffed00 0%, #f26529 75%); height: 2px !important; word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0; border: 0;" valign="top" align="left"> </td>
						</tr>
						<tr style="vertical-align: top; padding: 0;" align="left">
							<td class="mini-centered-text" style="color: #343b41; mso-table-lspace: 0pt; mso-table-rspace: 0pt; word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 25px 35px; font: 400 16px/27px 'Helvetica Neue', Helvetica, Arial, sans-serif;" align="center" valign="top">
								

\{\{Subject .Subject "\{\{.InvitedBy\}\} has invited you to join Grafana"\}\}

<table class="row" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; display: block; padding: 0px;">
	<tr style="vertical-align: top; padding: 0;" align="left">
		<td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">

			<table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
				<tr style="vertical-align: top; padding: 0;" align="left">
					<td style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="left" valign="top">
						<h4 class="center" style="color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 1.3; word-break: normal; font-size: 20px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="center">You're invited to join \{\{.OrgName\}\}</h4>
					</td>
					<td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
				</tr>
			</table>

		</td>
	</tr>
</table>

<table class="row" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; display: block; padding: 0px;">
	<tr style="vertical-align: top; padding: 0;" align="left">
		<td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">
			<table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
				<tr style="vertical-align: top; padding: 0;" align="left">
					<td class="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="center" valign="top">
						<p style="color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0 0 10px; padding: 0;" align="left">You've been invited to join the <b>\{\{.OrgName\}\}</b> organization by <b>\{\{.InvitedBy\}\}</b>. To accept your invitation and join the team, please click the link below:</p>
					</td>
				</tr>
				<tr style="vertical-align: top; padding: 0;" align="left">				
					<td class="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="center" valign="top">
	                    <table class="better-button" align="center" border="0" cellspacing="0" cellpadding="0" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; margin-top: 10px; margin-bottom: 20px; padding: 0;">
	                    	<tr style="vertical-align: top; padding: 0;" align="left">
	                          <td align="center" class="better-button" bgcolor="#ff8f2b" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; -webkit-border-radius: 2px; -moz-border-radius: 2px; border-radius: 2px; margin: 0; padding: 0px;" valign="top"><a href="\{\{.LinkUrl\}\}" target="_blank" style="color: #FFF; text-decoration: none; -webkit-border-radius: 2px; -moz-border-radius: 2px; border-radius: 2px; display: inline-block; padding: 12px 25px; border: 1px solid #ff8f2b;">Accept Invitation</a></td>
	                        </tr>
	                    </table>
					</td>
				</tr>
				<tr style="vertical-align: top; padding: 0;" align="left">
					<td class="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="center" valign="top">
						<p style="color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0 0 10px; padding: 0;" align="left">You can also copy/paste this link into your browser directly: <a href="\{\{.LinkUrl\}\}" style="color: #E67612; text-decoration: none;">\{\{.LinkUrl\}\}</a></p>
					</td>								
					<td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
				</tr>
			</table>
		</td>
	</tr>
</table>
								
							</td>
						</tr>
					</table>
					
					<table class="footer center" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: center; color: #999999; margin-top: 20px; padding: 0;" bgcolor="#2e2e2e">
						<tr style="vertical-align: top; padding: 0;" align="left">
							<td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 20px 0px 0px;" align="left" valign="top">
								<table class="twelve columns center" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: center; width: 580px; margin: 0 auto; padding: 0;">
									<tr style="vertical-align: top; padding: 0;" align="left">
										<td class="twelve" align="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" valign="top">
											<center style="width: 100%; min-width: 580px;">
												<p style="font-size: 12px; color: #999999; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0 0 10px; padding: 0;" align="center">
													Sent by <a href="\{\{.AppUrl\}\}" style="color: #E67612; text-decoration: none;">Grafana v\{\{.BuildVersion\}\}</a>
													<br />© 2016 Grafana and raintank
												</p>
											</center>
										</td>
										<td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
									</tr>
								</table>
							</td>
						</tr>
					</table>
				</center>
			</td>
		</tr>
	</table>
</body>
</html>

`
	PublicResetPassword = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns="http://www.w3.org/1999/xhtml">
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<meta name="viewport" content="width=device-width" />
	
<style>body {
width: 100% !important; min-width: 100%; -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; margin: 0; padding: 0;
}
img {
outline: none; text-decoration: none; -ms-interpolation-mode: bicubic; width: auto; float: left; clear: both; display: block;
}
body {
color: #222222; font-family: "Helvetica", "Arial", sans-serif; font-weight: normal; padding: 0; margin: 0; text-align: left; line-height: 1.3;
}
body {
font-size: 14px; line-height: 19px;
}
a:hover {
color: #2795b6 !important;
}
a:active {
color: #2795b6 !important;
}
a:visited {
color: #2ba6cb !important;
}
body {
font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none;
}
a:hover {
color: #ff8f2b !important;
}
a:active {
color: #F2821E !important;
}
a:visited {
color: #E67612 !important;
}
.better-button:hover a {
color: #FFFFFF !important; background-color: #F2821E; border: 1px solid #F2821E;
}
.better-button:visited a {
color: #FFFFFF !important;
}
.better-button:active a {
color: #FFFFFF !important;
}
.better-button-alt:hover a {
color: #ff8f2b !important; background-color: #DDDDDD; border: 1px solid #F2821E;
}
.better-button-alt:visited a {
color: #ff8f2b !important;
}
.better-button-alt:active a {
color: #ff8f2b !important;
}
body {
height: 100% !important; width: 100% !important;
}
body .copy {
-ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;
}
.ExternalClass {
width: 100%;
}
.ExternalClass {
line-height: 100%;
}
img {
-ms-interpolation-mode: bicubic;
}
img {
border: 0 !important; outline: none !important; text-decoration: none !important;
}
a:hover {
text-decoration: underline;
}
@media only screen and (max-width: 600px) {
  table[class="body"] center {
    min-width: 0 !important;
  }
  table[class="body"] .container {
    width: 95% !important;
  }
  table[class="body"] .row {
    width: 100% !important; display: block !important;
  }
  table[class="body"] .wrapper {
    display: block !important; padding-right: 0 !important;
  }
  table[class="body"] .columns {
    table-layout: fixed !important; float: none !important; width: 100% !important; padding-right: 0px !important; padding-left: 0px !important; display: block !important;
  }
  table[class="body"] table.columns td {
    width: 100% !important;
  }
  table[class="body"] .columns td.six {
    width: 50% !important;
  }
  table[class="body"] .columns td.twelve {
    width: 100% !important;
  }
  table[class="body"] table.columns td.expander {
    width: 1px !important;
  }
  .logo {
    margin-left: 10px;
  }
}
@media (max-width: 600px) {
  table[class="email-container"] {
    width: 95% !important;
  }
  img[class="fluid"] {
    width: 100% !important; max-width: 100% !important; height: auto !important; margin: auto !important;
  }
  img[class="fluid-centered"] {
    width: 100% !important; max-width: 100% !important; height: auto !important; margin: auto !important;
  }
  img[class="fluid-centered"] {
    margin: auto !important;
  }
  td[class="comms-content"] {
    padding: 20px !important;
  }
  td[class="stack-column"] {
    display: block !important; width: 100% !important; direction: ltr !important;
  }
  td[class="stack-column-center"] {
    display: block !important; width: 100% !important; direction: ltr !important;
  }
  td[class="stack-column-center"] {
    text-align: center !important;
  }
  td[class="copy"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="copy -center"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="copy -bold"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="small-text"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="mini-centered-text"] {
    font-size: 14px !important; line-height: 24px !important; padding: 15px 30px !important;
  }
  td[class="copy -padd"] {
    padding: 0 40px !important;
  }
  span[class="sep"] {
    display: none !important;
  }
  td[class="mb-hide"] {
    display: none !important; height: 0 !important;
  }
  td[class="spacer mb-shorten"] {
    height: 25px !important;
  }
  .two-up td {
    width: 270px;
  }
}
</style></head>
<body leftmargin="0" topmargin="0" marginwidth="0" marginheight="0" class="main" style="height: 100% !important; width: 100% !important; min-width: 100%; -webkit-text-size-adjust: none; -ms-text-size-adjust: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; text-align: left; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; margin: 0 auto; padding: 0;" bgcolor="#2e2e2e">

	<table class="body" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; height: 100%; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" bgcolor="#2e2e2e">
		<tr style="vertical-align: top; padding: 0;" align="left">
			<td class="center" align="center" valign="top" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;">
        <center style="width: 100%; min-width: 580px;">
					<table class="row header" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; margin-top: 25px; margin-bottom: 25px; padding: 0px;">
						<tr style="vertical-align: top; padding: 0;" align="left">
						  <td class="center" align="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" valign="top">
						    <center style="width: 100%; min-width: 580px;">

						      <table class="container" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: inherit; width: 580px; margin: 0 auto; padding: 0;">
						        <tr style="vertical-align: top; padding: 0;" align="left">
						          <td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">

						            <table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
						              <tr style="vertical-align: top; padding: 0;" align="left">
						                <td class="twelve sub-columns center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; min-width: 0px; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 10px 10px 0px;" align="center" valign="top">
                              <img class="logo" src="http://grafana.org/assets/img/logo_new_transparent_200x48.png" style="width: 200px; display: inline; outline: none !important; text-decoration: none !important; -ms-interpolation-mode: bicubic; clear: both; border: 0;" align="none" />
                            </td>
                            <td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
                          </tr>
						            </table>

						          </td>
						        </tr>
						      </table>

						    </center>
						  </td>
						</tr>
					</table>

					<table class="container" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: inherit; width: 580px; margin: 0 auto; padding: 0;" width="600" bgcolor="#efefef">
						<tr style="vertical-align: top; padding: 0;" align="left">
							<td height="2" class="spacer mb-shorten" style="font-size: 0; line-height: 0; mso-table-lspace: 0pt; mso-table-rspace: 0pt; background-image: linear-gradient(to right, #ffed00 0%, #f26529 75%); height: 2px !important; word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0; border: 0;" valign="top" align="left"> </td>
						</tr>
						<tr style="vertical-align: top; padding: 0;" align="left">
							<td class="mini-centered-text" style="color: #343b41; mso-table-lspace: 0pt; mso-table-rspace: 0pt; word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 25px 35px; font: 400 16px/27px 'Helvetica Neue', Helvetica, Arial, sans-serif;" align="center" valign="top">
								\{\{Subject .Subject "Reset your Grafana password - \{\{.Name\}\}"\}\}

<table class="row" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; display: block; padding: 0px;">
	<tr style="vertical-align: top; padding: 0;" align="left">
		<td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">

			<table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
				<tr style="vertical-align: top; padding: 0;" align="left">
					<td style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="left" valign="top">
						<h4 style="color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 1.3; word-break: normal; font-size: 20px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left">Hi \{\{.Name\}\},</h4>
					</td>
					<td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
				</tr>
			</table>

		</td>
	</tr>
</table>

<table class="row" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; display: block; padding: 0px;">
	<tr style="vertical-align: top; padding: 0;" align="left">
		<td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">
			<table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
				<tr style="vertical-align: top; padding: 0;" align="left">
					<td class="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="center" valign="top">
						<p style="color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0 0 10px; padding: 0;" align="left">
							Please click the following link to reset your password within <b>\{\{.EmailCodeValidHours\}\} hours</b>.
						</p>
						<p style="color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0 0 10px; padding: 0;" align="left">
							<a href="\{\{.AppUrl\}\}user/password/reset?code=\{\{.Code\}\}" style="color: #E67612; text-decoration: none;">\{\{.AppUrl\}\}user/password/reset?code=\{\{.Code\}\}</a>
						</p>
						<p style="color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0 0 10px; padding: 0;" align="left">Not working? Try copying and pasting it to your browser.</p>
					</td>
					<td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
				</tr>
			</table>

		</td>
	</tr>
</table>



								
							</td>
						</tr>
					</table>
					
					<table class="footer center" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: center; color: #999999; margin-top: 20px; padding: 0;" bgcolor="#2e2e2e">
						<tr style="vertical-align: top; padding: 0;" align="left">
							<td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 20px 0px 0px;" align="left" valign="top">
								<table class="twelve columns center" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: center; width: 580px; margin: 0 auto; padding: 0;">
									<tr style="vertical-align: top; padding: 0;" align="left">
										<td class="twelve" align="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" valign="top">
											<center style="width: 100%; min-width: 580px;">
												<p style="font-size: 12px; color: #999999; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0 0 10px; padding: 0;" align="center">
													Sent by <a href="\{\{.AppUrl\}\}" style="color: #E67612; text-decoration: none;">Grafana v\{\{.BuildVersion\}\}</a>
													<br />© 2016 Grafana and raintank
												</p>
											</center>
										</td>
										<td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
									</tr>
								</table>
							</td>
						</tr>
					</table>
				</center>
			</td>
		</tr>
	</table>
</body>
</html>

`
	PublicSignupStarted = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns="http://www.w3.org/1999/xhtml">
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<meta name="viewport" content="width=device-width" />
	
<style>body {
width: 100% !important; min-width: 100%; -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; margin: 0; padding: 0;
}
img {
outline: none; text-decoration: none; -ms-interpolation-mode: bicubic; width: auto; float: left; clear: both; display: block;
}
body {
color: #222222; font-family: "Helvetica", "Arial", sans-serif; font-weight: normal; padding: 0; margin: 0; text-align: left; line-height: 1.3;
}
body {
font-size: 14px; line-height: 19px;
}
a:hover {
color: #2795b6 !important;
}
a:active {
color: #2795b6 !important;
}
a:visited {
color: #2ba6cb !important;
}
body {
font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none;
}
a:hover {
color: #ff8f2b !important;
}
a:active {
color: #F2821E !important;
}
a:visited {
color: #E67612 !important;
}
.better-button:hover a {
color: #FFFFFF !important; background-color: #F2821E; border: 1px solid #F2821E;
}
.better-button:visited a {
color: #FFFFFF !important;
}
.better-button:active a {
color: #FFFFFF !important;
}
.better-button-alt:hover a {
color: #ff8f2b !important; background-color: #DDDDDD; border: 1px solid #F2821E;
}
.better-button-alt:visited a {
color: #ff8f2b !important;
}
.better-button-alt:active a {
color: #ff8f2b !important;
}
body {
height: 100% !important; width: 100% !important;
}
body .copy {
-ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;
}
.ExternalClass {
width: 100%;
}
.ExternalClass {
line-height: 100%;
}
img {
-ms-interpolation-mode: bicubic;
}
img {
border: 0 !important; outline: none !important; text-decoration: none !important;
}
a:hover {
text-decoration: underline;
}
@media only screen and (max-width: 600px) {
  table[class="body"] center {
    min-width: 0 !important;
  }
  table[class="body"] .container {
    width: 95% !important;
  }
  table[class="body"] .row {
    width: 100% !important; display: block !important;
  }
  table[class="body"] .wrapper {
    display: block !important; padding-right: 0 !important;
  }
  table[class="body"] .columns {
    table-layout: fixed !important; float: none !important; width: 100% !important; padding-right: 0px !important; padding-left: 0px !important; display: block !important;
  }
  table[class="body"] table.columns td {
    width: 100% !important;
  }
  table[class="body"] .columns td.six {
    width: 50% !important;
  }
  table[class="body"] .columns td.twelve {
    width: 100% !important;
  }
  table[class="body"] table.columns td.expander {
    width: 1px !important;
  }
  .logo {
    margin-left: 10px;
  }
}
@media (max-width: 600px) {
  table[class="email-container"] {
    width: 95% !important;
  }
  img[class="fluid"] {
    width: 100% !important; max-width: 100% !important; height: auto !important; margin: auto !important;
  }
  img[class="fluid-centered"] {
    width: 100% !important; max-width: 100% !important; height: auto !important; margin: auto !important;
  }
  img[class="fluid-centered"] {
    margin: auto !important;
  }
  td[class="comms-content"] {
    padding: 20px !important;
  }
  td[class="stack-column"] {
    display: block !important; width: 100% !important; direction: ltr !important;
  }
  td[class="stack-column-center"] {
    display: block !important; width: 100% !important; direction: ltr !important;
  }
  td[class="stack-column-center"] {
    text-align: center !important;
  }
  td[class="copy"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="copy -center"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="copy -bold"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="small-text"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="mini-centered-text"] {
    font-size: 14px !important; line-height: 24px !important; padding: 15px 30px !important;
  }
  td[class="copy -padd"] {
    padding: 0 40px !important;
  }
  span[class="sep"] {
    display: none !important;
  }
  td[class="mb-hide"] {
    display: none !important; height: 0 !important;
  }
  td[class="spacer mb-shorten"] {
    height: 25px !important;
  }
  .two-up td {
    width: 270px;
  }
}
</style></head>
<body leftmargin="0" topmargin="0" marginwidth="0" marginheight="0" class="main" style="height: 100% !important; width: 100% !important; min-width: 100%; -webkit-text-size-adjust: none; -ms-text-size-adjust: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; text-align: left; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; margin: 0 auto; padding: 0;" bgcolor="#2e2e2e">

	<table class="body" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; height: 100%; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" bgcolor="#2e2e2e">
		<tr style="vertical-align: top; padding: 0;" align="left">
			<td class="center" align="center" valign="top" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;">
        <center style="width: 100%; min-width: 580px;">
					<table class="row header" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; margin-top: 25px; margin-bottom: 25px; padding: 0px;">
						<tr style="vertical-align: top; padding: 0;" align="left">
						  <td class="center" align="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" valign="top">
						    <center style="width: 100%; min-width: 580px;">

						      <table class="container" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: inherit; width: 580px; margin: 0 auto; padding: 0;">
						        <tr style="vertical-align: top; padding: 0;" align="left">
						          <td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">

						            <table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
						              <tr style="vertical-align: top; padding: 0;" align="left">
						                <td class="twelve sub-columns center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; min-width: 0px; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 10px 10px 0px;" align="center" valign="top">
                              <img class="logo" src="http://grafana.org/assets/img/logo_new_transparent_200x48.png" style="width: 200px; display: inline; outline: none !important; text-decoration: none !important; -ms-interpolation-mode: bicubic; clear: both; border: 0;" align="none" />
                            </td>
                            <td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
                          </tr>
						            </table>

						          </td>
						        </tr>
						      </table>

						    </center>
						  </td>
						</tr>
					</table>

					<table class="container" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: inherit; width: 580px; margin: 0 auto; padding: 0;" width="600" bgcolor="#efefef">
						<tr style="vertical-align: top; padding: 0;" align="left">
							<td height="2" class="spacer mb-shorten" style="font-size: 0; line-height: 0; mso-table-lspace: 0pt; mso-table-rspace: 0pt; background-image: linear-gradient(to right, #ffed00 0%, #f26529 75%); height: 2px !important; word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0; border: 0;" valign="top" align="left"> </td>
						</tr>
						<tr style="vertical-align: top; padding: 0;" align="left">
							<td class="mini-centered-text" style="color: #343b41; mso-table-lspace: 0pt; mso-table-rspace: 0pt; word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 25px 35px; font: 400 16px/27px 'Helvetica Neue', Helvetica, Arial, sans-serif;" align="center" valign="top">
								\{\{Subject .Subject "Welcome to Grafana, please complete your sign up!"\}\}

<table class="row" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; display: block; padding: 0px;">
	<tr style="vertical-align: top; padding: 0;" align="left">
		<td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">

			<table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
				<tr style="vertical-align: top; padding: 0;" align="left">
					<td style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="left" valign="top">
						<h4 class="center" style="color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 1.3; word-break: normal; font-size: 20px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="center">Complete the signup</h4>
					</td>
					<td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
				</tr>
			</table>

		</td>
	</tr>
</table>

<table class="row" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; display: block; padding: 0px;">
	<tr style="vertical-align: top; padding: 0;" align="left">
		<td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">
			<table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
				<tr style="vertical-align: top; padding: 0;" align="left">
					<td class="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="center" valign="top">
						Copy and past the email verification code:<br />
						<span class="verification-code" style="background-color: #EEEEEE; display: inline-block; font-weight: bold; font-size: 20px; margin: 8px; padding: 3px;">\{\{.Code\}\}</span><br /> in
						the sign up form <strong>or</strong> use the link below.
					</td>
					<td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
				</tr>
				<tr style="vertical-align: top; padding: 0;" align="left">
					<td class="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="center" valign="top">
						<table class="better-button" align="center" border="0" cellspacing="0" cellpadding="0" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; margin-top: 10px; margin-bottom: 20px; padding: 0;">
							<tr style="vertical-align: top; padding: 0;" align="left">
								<td align="center" class="better-button" bgcolor="#ff8f2b" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; -webkit-border-radius: 2px; -moz-border-radius: 2px; border-radius: 2px; margin: 0; padding: 0px;" valign="top"><a href="\{\{.SignUpUrl\}\}" target="_blank" style="color: #FFF; text-decoration: none; -webkit-border-radius: 2px; -moz-border-radius: 2px; border-radius: 2px; display: inline-block; padding: 12px 25px; border: 1px solid #ff8f2b;">Complete Sign Up</a></td>
							</tr>
						</table>
					</td>
				</tr>
			</table>
		</td>
	</tr>
</table>



								
							</td>
						</tr>
					</table>
					
					<table class="footer center" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: center; color: #999999; margin-top: 20px; padding: 0;" bgcolor="#2e2e2e">
						<tr style="vertical-align: top; padding: 0;" align="left">
							<td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 20px 0px 0px;" align="left" valign="top">
								<table class="twelve columns center" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: center; width: 580px; margin: 0 auto; padding: 0;">
									<tr style="vertical-align: top; padding: 0;" align="left">
										<td class="twelve" align="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" valign="top">
											<center style="width: 100%; min-width: 580px;">
												<p style="font-size: 12px; color: #999999; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0 0 10px; padding: 0;" align="center">
													Sent by <a href="\{\{.AppUrl\}\}" style="color: #E67612; text-decoration: none;">Grafana v\{\{.BuildVersion\}\}</a>
													<br />© 2016 Grafana and raintank
												</p>
											</center>
										</td>
										<td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
									</tr>
								</table>
							</td>
						</tr>
					</table>
				</center>
			</td>
		</tr>
	</table>
</body>
</html>

`
	PublicWelcomeOnSignup = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns="http://www.w3.org/1999/xhtml">
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<meta name="viewport" content="width=device-width" />
	
<style>body {
width: 100% !important; min-width: 100%; -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; margin: 0; padding: 0;
}
img {
outline: none; text-decoration: none; -ms-interpolation-mode: bicubic; width: auto; float: left; clear: both; display: block;
}
body {
color: #222222; font-family: "Helvetica", "Arial", sans-serif; font-weight: normal; padding: 0; margin: 0; text-align: left; line-height: 1.3;
}
body {
font-size: 14px; line-height: 19px;
}
a:hover {
color: #2795b6 !important;
}
a:active {
color: #2795b6 !important;
}
a:visited {
color: #2ba6cb !important;
}
body {
font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none;
}
a:hover {
color: #ff8f2b !important;
}
a:active {
color: #F2821E !important;
}
a:visited {
color: #E67612 !important;
}
.better-button:hover a {
color: #FFFFFF !important; background-color: #F2821E; border: 1px solid #F2821E;
}
.better-button:visited a {
color: #FFFFFF !important;
}
.better-button:active a {
color: #FFFFFF !important;
}
.better-button-alt:hover a {
color: #ff8f2b !important; background-color: #DDDDDD; border: 1px solid #F2821E;
}
.better-button-alt:visited a {
color: #ff8f2b !important;
}
.better-button-alt:active a {
color: #ff8f2b !important;
}
body {
height: 100% !important; width: 100% !important;
}
body .copy {
-ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;
}
.ExternalClass {
width: 100%;
}
.ExternalClass {
line-height: 100%;
}
img {
-ms-interpolation-mode: bicubic;
}
img {
border: 0 !important; outline: none !important; text-decoration: none !important;
}
a:hover {
text-decoration: underline;
}
@media only screen and (max-width: 600px) {
  table[class="body"] center {
    min-width: 0 !important;
  }
  table[class="body"] .container {
    width: 95% !important;
  }
  table[class="body"] .row {
    width: 100% !important; display: block !important;
  }
  table[class="body"] .wrapper {
    display: block !important; padding-right: 0 !important;
  }
  table[class="body"] .columns {
    table-layout: fixed !important; float: none !important; width: 100% !important; padding-right: 0px !important; padding-left: 0px !important; display: block !important;
  }
  table[class="body"] table.columns td {
    width: 100% !important;
  }
  table[class="body"] .columns td.six {
    width: 50% !important;
  }
  table[class="body"] .columns td.twelve {
    width: 100% !important;
  }
  table[class="body"] table.columns td.expander {
    width: 1px !important;
  }
  .logo {
    margin-left: 10px;
  }
}
@media (max-width: 600px) {
  table[class="email-container"] {
    width: 95% !important;
  }
  img[class="fluid"] {
    width: 100% !important; max-width: 100% !important; height: auto !important; margin: auto !important;
  }
  img[class="fluid-centered"] {
    width: 100% !important; max-width: 100% !important; height: auto !important; margin: auto !important;
  }
  img[class="fluid-centered"] {
    margin: auto !important;
  }
  td[class="comms-content"] {
    padding: 20px !important;
  }
  td[class="stack-column"] {
    display: block !important; width: 100% !important; direction: ltr !important;
  }
  td[class="stack-column-center"] {
    display: block !important; width: 100% !important; direction: ltr !important;
  }
  td[class="stack-column-center"] {
    text-align: center !important;
  }
  td[class="copy"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="copy -center"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="copy -bold"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="small-text"] {
    font-size: 14px !important; line-height: 24px !important; padding: 0 30px !important;
  }
  td[class="mini-centered-text"] {
    font-size: 14px !important; line-height: 24px !important; padding: 15px 30px !important;
  }
  td[class="copy -padd"] {
    padding: 0 40px !important;
  }
  span[class="sep"] {
    display: none !important;
  }
  td[class="mb-hide"] {
    display: none !important; height: 0 !important;
  }
  td[class="spacer mb-shorten"] {
    height: 25px !important;
  }
  .two-up td {
    width: 270px;
  }
}
</style></head>
<body leftmargin="0" topmargin="0" marginwidth="0" marginheight="0" class="main" style="height: 100% !important; width: 100% !important; min-width: 100%; -webkit-text-size-adjust: none; -ms-text-size-adjust: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; text-align: left; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; margin: 0 auto; padding: 0;" bgcolor="#2e2e2e">

	<table class="body" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; height: 100%; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" bgcolor="#2e2e2e">
		<tr style="vertical-align: top; padding: 0;" align="left">
			<td class="center" align="center" valign="top" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;">
        <center style="width: 100%; min-width: 580px;">
					<table class="row header" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; margin-top: 25px; margin-bottom: 25px; padding: 0px;">
						<tr style="vertical-align: top; padding: 0;" align="left">
						  <td class="center" align="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" valign="top">
						    <center style="width: 100%; min-width: 580px;">

						      <table class="container" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: inherit; width: 580px; margin: 0 auto; padding: 0;">
						        <tr style="vertical-align: top; padding: 0;" align="left">
						          <td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">

						            <table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
						              <tr style="vertical-align: top; padding: 0;" align="left">
						                <td class="twelve sub-columns center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; min-width: 0px; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 10px 10px 0px;" align="center" valign="top">
                              <img class="logo" src="http://grafana.org/assets/img/logo_new_transparent_200x48.png" style="width: 200px; display: inline; outline: none !important; text-decoration: none !important; -ms-interpolation-mode: bicubic; clear: both; border: 0;" align="none" />
                            </td>
                            <td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
                          </tr>
						            </table>

						          </td>
						        </tr>
						      </table>

						    </center>
						  </td>
						</tr>
					</table>

					<table class="container" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: inherit; width: 580px; margin: 0 auto; padding: 0;" width="600" bgcolor="#efefef">
						<tr style="vertical-align: top; padding: 0;" align="left">
							<td height="2" class="spacer mb-shorten" style="font-size: 0; line-height: 0; mso-table-lspace: 0pt; mso-table-rspace: 0pt; background-image: linear-gradient(to right, #ffed00 0%, #f26529 75%); height: 2px !important; word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0; border: 0;" valign="top" align="left"> </td>
						</tr>
						<tr style="vertical-align: top; padding: 0;" align="left">
							<td class="mini-centered-text" style="color: #343b41; mso-table-lspace: 0pt; mso-table-rspace: 0pt; word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 25px 35px; font: 400 16px/27px 'Helvetica Neue', Helvetica, Arial, sans-serif;" align="center" valign="top">
								\{\{Subject .Subject "Welcome to Grafana"\}\}

<table class="row" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; display: block; padding: 0px;">
	<tr style="vertical-align: top; padding: 0;" align="left">
		<td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">

			<table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
				<tr style="vertical-align: top; padding: 0;" align="left">
					<td style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="left" valign="top">
						<h4 style="color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 1.3; word-break: normal; font-size: 20px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left">Hi \{\{.Name\}\},</h4>
					</td>
					<td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
				</tr>
				<tr style="vertical-align: top; padding: 0;" align="left">
					<td style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="left" valign="top">
						Welcome! Ready to start building some beautiful metric and analytic dashboards?
					</td>
				</tr>
			</table>

		</td>
	</tr>
</table>

<table class="row" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 100%; position: relative; display: block; padding: 0px;">
	<tr style="vertical-align: top; padding: 0;" align="left">
		<td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 0px 0px;" align="left" valign="top">
			<table class="twelve columns" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: left; width: 580px; margin: 0 auto; padding: 0;">
				<tr style="vertical-align: top; padding: 0;" align="left">
					<td class="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="center" valign="top">
						<p style="color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0 0 10px; padding: 0;" align="left">
							If you are new to Grafana please read the <a href="http://docs.grafana.org/guides/gettingstarted/" style="color: #E67612; text-decoration: none;">Getting Started</a> guide.
						</p>
					</td>
					<td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
				</tr>
				<tr style="vertical-align: top; padding: 0;" align="left">
					<td style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" align="left" valign="top">
						Thank you for joining our community.
						<br />
						<p style="color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0 0 10px; padding: 0;" align="left">The Grafana Team</p>
					</td>
				</tr>
			</table>
		</td>
	</tr>
</table>



								
							</td>
						</tr>
					</table>
					
					<table class="footer center" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: center; color: #999999; margin-top: 20px; padding: 0;" bgcolor="#2e2e2e">
						<tr style="vertical-align: top; padding: 0;" align="left">
							<td class="wrapper last" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; position: relative; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 10px 20px 0px 0px;" align="left" valign="top">
								<table class="twelve columns center" style="border-spacing: 0; border-collapse: collapse; vertical-align: top; text-align: center; width: 580px; margin: 0 auto; padding: 0;">
									<tr style="vertical-align: top; padding: 0;" align="left">
										<td class="twelve" align="center" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; width: 100%; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0px 0px 10px;" valign="top">
											<center style="width: 100%; min-width: 580px;">
												<p style="font-size: 12px; color: #999999; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0 0 10px; padding: 0;" align="center">
													Sent by <a href="\{\{.AppUrl\}\}" style="color: #E67612; text-decoration: none;">Grafana v\{\{.BuildVersion\}\}</a>
													<br />© 2016 Grafana and raintank
												</p>
											</center>
										</td>
										<td class="expander" style="word-break: break-word; -webkit-hyphens: auto; -moz-hyphens: auto; hyphens: auto; border-collapse: collapse !important; visibility: hidden; width: 0px; color: #222222; font-family: 'Open Sans', 'Helvetica Neue', 'Helvetica', Helvetica, Arial, sans-serif; font-weight: normal; line-height: 19px; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; margin: 0; padding: 0;" align="left" valign="top"></td>
									</tr>
								</table>
							</td>
						</tr>
					</table>
				</center>
			</td>
		</tr>
	</table>
</body>
</html>

`
)
