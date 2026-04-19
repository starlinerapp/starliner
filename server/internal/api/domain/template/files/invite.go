package files

const InviteEmail = `<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 40px 20px;">
  <h1 style="color: #333; font-size: 24px;">You're invited!</h1>
  <p style="color: #555; font-size: 16px; line-height: 1.5;">
    You've been invited to join <strong>{{.OrganizationName}}</strong> on Starliner.
  </p>
  <a href="{{.InviteLink}}" style="display: inline-block; background-color: #6955C4; color: #fff; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin-top: 20px;">
    Accept Invite
  </a>
  <p style="color: #999; font-size: 12px; margin-top: 40px;">
    If you didn't expect this invitation, you can ignore this email.
  </p>
</div>`
