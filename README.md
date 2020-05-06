# email-notification program using Gmail's API

Often you might want, or need to, send some kind of notification email from
your server. My current usecase is that I manage a VPS, and I want to receive
certain notification emails from time to time after certain actions happen.

There's multiple ways to go about this. You can use SMTP and log into a
supported email account you own. But that requires saving the credentials for
the whole password somewhere.

You can set up some MTA in your server and configure it (not for the faint of
heart). It also requires installing the software, which might be heavy, and
then depend on your server's IP not being into any email blacklists, so that
your emails will actually get delivered.

Doing some research I found out Google has APIs for most of their services,
and Gmail is one of those. It's kind of a pain to set up, but it works pretty
well once you save the OAuth hurdles.

## Setting up

Visit the [Google Developers Console][]. You will have to create a project.
Then, in the Credentials tab set up new OAuth2 credentials for an `Other` type
of app.

You might need to fill something about "OAuth2 Consent" before you're allowed
to set up new OAuth2 credentials. As far as I can tell, it's to publish your
"app", but you can just go through with the first step that asks you to fill
fields like the app name, and then **not** go through with the
verification/publication step.

Once you have the OAuth2 credentials for the `Other` app, go into it's edit
page, and download the credentials JSON file. Store that file anywhere, and
note the path down.

## Usage

This makes use of OAuth2 so it requires an initial set up. First the
credentials JSON file generated from the Developers interface. This contains
OAuth params such as the Client ID and Secret ID needed for the generation of
the OAuth2 tokens the program will then use to send the emails.

I'd suggest storing this file as `~/.email-notification/credentials.json`. But
the path to fetch this file can nevertheless be configured through the
`--credentials` flag.

First usage will require generating the `token.json` file, which contains
authentication tokens. This can be done at any computer, but I'd suggest doing
it at your local computer, and then uploading the binary, `credentials.json`,
and the resulting `token.json` to the server.

The authentication will require opening a given URL that will show up in the
console in a browser, and then copying the string shown in that browser tab
and pastying it back into the console (will be waiting for input unless you've
done anything else). Paste it in, hit Enter, and it should have gone ahead
with the token generation and creating the `token.json` file.

The default lookup path for the `token.json` file is also under
`~/.email-notification/`, so I'd suggest storing it there as well. But the
lookup path can be changed through the `--token` flag.

## Install
### Using Go

```bash
git clone https://github.com/Sighery/email-notification.git
cd email-notification
make build
# Now you can use the binary
./email-notification --help
# Or move it to /usr/bin to use it system-wide
cp email-notification /usr/bin/.
email-notification --help
```

### Using Arch Linux

```bash
git clone https://github.com/Sighery/email-notification.git
cd email-notification
makepkg --needed --syncdeps --install
email-notification --help
```

Alternatively, it'll also be up on my [own Arch Linux repository][Arch
Linux repository].

[Google Developers Console]: https://console.developers.google.com/
[Arch Linux repository]: https://archrepo.sighery.com/
