# npoleon
Npoleon is a command-line utility that helps you scrobble tracks that are
being, have been, or will be played on NPO Radio 1, 2 and 3FM to Last.fm, a
social network centred around music that – like 3FM – inexplicably still exists
despite more than a decade of declining market share.

This tool was written over the course of several train rides in the Netherlands,
mainly as a fun way to learn more about Go and to find out if I still dislike
the language. Sadly it turns out I do, so don’t expect (m)any updates unless
something is completely broken. Although Npoleon has a few rough edges here and
there, I’ve been [dogfooding] Npoleon for a while now so you probably shouldn’t
run into any major issues.

[dogfooding]: https://en.wikipedia.org/wiki/Eating_your_own_dog_food

## Installation
It only takes three steps to set up Npoleon for scrobbling.

First, download and install the latest version of Npoleon by downloading the
right binary for your platform. Ensure that the binary is somewhere in your
path (e.g. `/usr/local/bin` and has the executable bit set (using `chmod +x`).

You will need to [create an API account][lastfm] in order to access Last.fm’s
Track API.  It doesn’t really matter what you enter here, because you will be
the tool’s only user.

[lastfm]: https://www.last.fm/api/account/create

Once you have created an API account, Last.fm will show you a key and a secret.
Write these to `~/.npoleon/config` as follows, making sure to replace the dummy
values with your actual key and secret):

```
LASTFM_API_KEY=0123456789abcdef0123456789abcdef
LASTFM_API_SECRET=0123456789abcdef0123456789abcdef
```

Finally, run this command and follow the instructions:

```
npoleon login
```

This gives your specific instance of Npoleon access to your Last.fm account
and allows it to scrobble tracks on your behalf. You only need to do this
once.

## Usage
Npoleon can scrobble tracks for three NPO radio stations: `nporadio1`,
`nporadio2`, and `npo3fm` (or `radio1`, `radio2` and `3fm`). The examples
below assume that you want to scrobble tracks for `3fm`.

To scrobble a single track that’s currently being played, execute:

```
npoleon scrobble 3fm --once
```

To keep scrobbling tracks indefinitely (at least until you terminate the
command), simply execute:

```
npoleon scrobble 3fm
```

You can also scrobble all tracks that have been played since a particular
moment:

```
npoleon scrobble 3fm --from "2024-01-20 14:30:00"
```

Or ask Npoleon to scrobble tracks until a specific time:

```
npoleon scrobble 3fm --until "2024-01-20 20:55:00"
```

`--from` and `--until` can be combined to scrobble tracks for specific periods:

```
npoleon scrobble 3fm --from "2024-01-20 14:30:00" --until "2024-01-20 20:55:00"
```
