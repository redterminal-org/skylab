# Hi Valeria!

Anyway, I think the gemmail client should work as follows. If you have ideas for how things should work differently let me know!

## Name: **Skylab** Misfin Client
gemini://satch.xyz/skylab/

## Front Page + client certs
- Url is `gemini://mail.example.org` or whatever, and the server serves an index file which is just some static gemtext if you do not have a client certificate activated
- If you click on a login link `gemini://mail.example.org/?register`, it prompts you to activate a client certificate using a 60 response (see similar behavior on [BBS](gemini://bbs.geminispace.org))
- Then if you activate a client certificate it does one of several things: 1) it gives you a 61 error message if you don't have an account on the server and signups are closed, or 2) it prompts you to create an account if signups or open, or 3) it shows you your inbox/home page thing if you have an account.
- The rest of the pages all redirect you to the homepage (30) if you have no client cert or give an error (61) if you don't have one at all

## Inbox layout and design
- Imitates tradition webmail interface. Front page should have links to Trash, Archive, Sent, and Drafts folder, as well as maybe 'More' or some type of tags system. Below that (hopefully not too far down) is links to read various new mail. Emoji icons can indicate which mail is unread, as well as built in gemini client read indicators such as are in Lagrange.
- Draft composer similar to BBS
