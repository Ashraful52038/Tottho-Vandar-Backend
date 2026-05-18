package email

const (
	CommentEmailTemplate = `<!DOCTYPE html>
<html lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="x-apple-disable-message-reformatting">
    <title>New Comment</title>
    <!--[if mso]>
    <noscript><xml><o:OfficeDocumentSettings><o:PixelsPerInch>96</o:PixelsPerInch></o:OfficeDocumentSettings></xml></noscript>
    <![endif]-->
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Playfair+Display:wght@600;700&family=DM+Sans:wght@300;400;500;600&display=swap');
        * { margin: 0; padding: 0; }
        body { background-color: #0b0f1a !important; }
        @media only screen and (max-width:600px) {
            .email-container { width: 100%% !important; }
            .pad-mobile { padding: 28px 20px !important; }
            .hero-pad { padding: 40px 20px 32px !important; }
        }
    </style>
</head>
<body style="margin:0;padding:0;background-color:#0b0f1a;">

<table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0" bgcolor="#0b0f1a">
<tr><td align="center" style="padding:48px 16px;">

  <table role="presentation" class="email-container" width="580" cellpadding="0" cellspacing="0" border="0" style="max-width:580px;width:100%%;">

    <!-- WORDMARK -->
    <tr>
      <td align="center" style="padding-bottom:26px;">
        <table role="presentation" cellpadding="0" cellspacing="0" border="0">
          <tr>
            <td style="padding-bottom:7px;border-bottom:1px solid #7a5c1e;">
              <span style="font-family:'Playfair Display',Georgia,serif;font-size:11px;font-weight:700;letter-spacing:6px;text-transform:uppercase;color:#c9a84c;">TOTTHO VANDAR</span>
            </td>
          </tr>
        </table>
      </td>
    </tr>

    <!-- MAIN CARD -->
    <tr>
      <td style="background-color:#111827;border-radius:6px;overflow:hidden;border:1px solid #1f2a3c;">

        <!-- HERO -->
        <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0">
          <tr>
            <td class="hero-pad" align="center" style="background-color:#0d1929;padding:52px 40px 44px;border-bottom:1px solid #1a2744;">

              <!-- Top accent line -->
              <table role="presentation" width="40" cellpadding="0" cellspacing="0" border="0" style="margin:0 auto 30px;">
                <tr><td height="2" bgcolor="#c9a84c" style="background-color:#c9a84c;font-size:0;line-height:0;">&nbsp;</td></tr>
              </table>

              <!-- Icon -->
              <table role="presentation" cellpadding="0" cellspacing="0" border="0" style="margin:0 auto 24px;">
                <tr>
                  <td align="center" valign="middle" style="background-color:#12243d;border:1px solid #1e3a5f;border-radius:6px;width:66px;height:66px;">
                    <p style="margin:0;font-size:28px;line-height:66px;width:66px;text-align:center;">&#128172;</p>
                  </td>
                </tr>
              </table>

              <p style="margin:0 0 10px;font-family:'DM Sans',Arial,sans-serif;font-size:10px;font-weight:600;letter-spacing:4px;text-transform:uppercase;color:#c9a84c;">NEW ACTIVITY</p>
              <h1 style="margin:0 0 10px;font-family:'Playfair Display',Georgia,serif;font-size:28px;font-weight:700;color:#e8eef6;line-height:1.2;">Someone left a comment</h1>
              <p style="margin:0;font-family:'DM Sans',Arial,sans-serif;font-size:13px;color:#3d5068;font-weight:300;letter-spacing:0.3px;">Your post sparked a conversation.</p>

            </td>
          </tr>
        </table>

        <!-- BODY -->
        <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0">
          <tr>
            <td class="pad-mobile" style="padding:36px 40px;">

              <!-- Section label: From -->
              <p style="margin:0 0 10px;font-family:'DM Sans',Arial,sans-serif;font-size:9px;font-weight:600;letter-spacing:3px;text-transform:uppercase;color:#2d3f52;">FROM</p>

              <!-- Who commented -->
              <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color:#0d1929;border:1px solid #1a2744;border-radius:4px;margin-bottom:30px;">
                <tr>
                  <td style="padding:18px 20px;">
                    <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0">
                      <tr>
                        <td width="46" valign="middle">
                          <table role="presentation" cellpadding="0" cellspacing="0" border="0">
                            <tr>
                              <td align="center" valign="middle" style="background-color:#12243d;border:1px solid #1e3a5f;border-radius:4px;width:46px;height:46px;">
                                <p style="margin:0;font-size:20px;line-height:46px;width:46px;text-align:center;">&#128100;</p>
                              </td>
                            </tr>
                          </table>
                        </td>
                        <td valign="middle" style="padding-left:16px;">
                          <p style="margin:0 0 4px;font-family:'Playfair Display',Georgia,serif;font-size:16px;font-weight:600;color:#cbd5e1;">%s</p>
                          <p style="margin:0;font-family:'DM Sans',Arial,sans-serif;font-size:10px;font-weight:600;letter-spacing:1.5px;text-transform:uppercase;color:#2d3f52;">Commented on your post</p>
                        </td>
                        <td align="right" valign="middle">
                          <table role="presentation" cellpadding="0" cellspacing="0" border="0">
                            <tr>
                              <td style="background-color:#1a1000;border:1px solid #7a5c1e;border-radius:20px;padding:6px 14px;">
                                <span style="font-family:'DM Sans',Arial,sans-serif;font-size:11px;font-weight:600;color:#c9a84c;letter-spacing:0.5px;">&#128172;&nbsp; Comment</span>
                              </td>
                            </tr>
                          </table>
                        </td>
                      </tr>
                    </table>
                  </td>
                </tr>
              </table>

              <!-- Section label: Post -->
              <p style="margin:0 0 10px;font-family:'DM Sans',Arial,sans-serif;font-size:9px;font-weight:600;letter-spacing:3px;text-transform:uppercase;color:#2d3f52;">ON YOUR POST</p>

              <!-- Post title -->
              <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color:#0d1929;border:1px solid #1a2744;border-radius:4px;margin-bottom:16px;overflow:hidden;">
                <tr>
                  <td width="3" bgcolor="#c9a84c" style="background-color:#c9a84c;font-size:0;line-height:0;">&nbsp;</td>
                  <td style="padding:18px 20px;">
                    <p style="margin:0;font-family:'Playfair Display',Georgia,serif;font-size:16px;font-weight:600;color:#94a3b8;line-height:1.35;">%s</p>
                  </td>
                </tr>
              </table>

              <!-- Section label: Comment -->
              <p style="margin:0 0 10px;font-family:'DM Sans',Arial,sans-serif;font-size:9px;font-weight:600;letter-spacing:3px;text-transform:uppercase;color:#2d3f52;">THEIR COMMENT</p>

              <!-- Comment bubble -->
              <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color:#0d1929;border:1px solid #1a2744;border-radius:4px;margin-bottom:34px;overflow:hidden;">
                <tr>
                  <td width="3" bgcolor="#c9a84c" style="background-color:#c9a84c;font-size:0;line-height:0;">&nbsp;</td>
                  <td style="padding:22px 22px 22px 20px;">
                    <!-- Opening quote mark -->
                    <p style="margin:0 0 10px;font-family:'Playfair Display',Georgia,serif;font-size:32px;line-height:1;color:#7a5c1e;">&ldquo;</p>
                    <p style="margin:0;font-family:'DM Sans',Arial,sans-serif;font-size:14px;color:#64748b;line-height:1.75;font-weight:300;font-style:italic;">%s</p>
                  </td>
                </tr>
              </table>

              <!-- CTA -->
              <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:32px;">
                <tr>
                  <td align="center">
                    <table role="presentation" cellpadding="0" cellspacing="0" border="0">
                      <tr>
                        <td style="background-color:#c9a84c;border-radius:2px;">
                          <a href="%s" style="display:inline-block;font-family:'DM Sans',Arial,sans-serif;font-size:12px;font-weight:700;color:#0b0f1a;text-decoration:none;padding:15px 44px;letter-spacing:2px;text-transform:uppercase;">View &amp; Reply &nbsp;&#8594;</a>
                        </td>
                      </tr>
                    </table>
                  </td>
                </tr>
              </table>

              <!-- Divider -->
              <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:22px;">
                <tr><td height="1" bgcolor="#1a2744" style="background-color:#1a2744;font-size:0;line-height:0;">&nbsp;</td></tr>
              </table>

              <!-- Footer note -->
              <p style="margin:0;font-family:'DM Sans',Arial,sans-serif;font-size:11px;color:#2d3f52;text-align:center;line-height:2;">
                You&#39;re receiving this because someone commented on your post on Tottho Vandar.<br>
                <a href="http://localhost:3000/settings" style="color:#7a5c1e;text-decoration:none;">Manage notifications</a>
                &nbsp;&nbsp;&#183;&nbsp;&nbsp;
                <a href="http://localhost:3000" style="color:#7a5c1e;text-decoration:none;">Visit Tottho Vandar</a>
              </p>

            </td>
          </tr>
        </table>

        <!-- CARD FOOTER -->
        <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0">
          <tr>
            <td align="center" style="background-color:#0a0e18;padding:26px 32px;border-top:1px solid #1a2744;">
              <table role="presentation" width="28" cellpadding="0" cellspacing="0" border="0" style="margin:0 auto 14px;">
                <tr><td height="1" bgcolor="#7a5c1e" style="background-color:#7a5c1e;font-size:0;line-height:0;">&nbsp;</td></tr>
              </table>
              <p style="margin:0 0 6px;font-family:'Playfair Display',Georgia,serif;font-size:12px;font-weight:700;color:#c9a84c;letter-spacing:4px;text-transform:uppercase;">Tottho Vandar</p>
              <p style="margin:0;font-family:'DM Sans',Arial,sans-serif;font-size:10px;color:#1e2d3d;letter-spacing:0.5px;">&#169; 2026 &nbsp;&#183;&nbsp; Treasure of Information</p>
            </td>
          </tr>
        </table>

      </td>
    </tr>

  </table>
</td></tr>
</table>

</body>
</html>`
)
