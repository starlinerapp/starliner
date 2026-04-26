package files

const InviteEmail = `<div style="font-family: Arial, sans-serif; max-width: 600px">
  <img src="https://starliner-596451156994-eu-north-1-an.s3.eu-north-1.amazonaws.com/starliner-logo.png" height="35px" style="margin-bottom: 20px"/>
  <h1 style="color: #202020; font-size: 16px;">Join {{.OrganizationName}} on Starliner</h1>
  <p style="color: #202020; font-size: 16px; line-height: 1.2;">
    You have been invited to collaborate in the organization {{.OrganizationName}} on Starliner.
  </p>
  <a href="{{.InviteLink}}" style="display: inline-block; background-color: #6955C4; color: #fff; padding: 12px 24px; text-decoration: none; border-radius: 6px;">
    Join this organization
  </a>
  <p style="color: #202020; font-size: 16px; margin-top: 20px;">
    If you didn't expect this invitation, you can safely ignore this email.
  </p>
  <p style="color: #202020; font-size: 16px; margin-top: 20px;">
    Cheers, <br />
    Starliner
  </p>
</div>`
