Configuration File
==================

DTrack has defaults that are designed to function adequately on a Raspberry
Pi 5, although these make very poor monitoring devices. If your [Selected
Hardware](hardware.md) provides a hardware video **en**conder, then it is
important to modify **record_save_options** (below).


config.json
-----------

Configuration is stored in a **JSON** file, called **config.json**.

**Important JSON Tips:**

1\. Format:
```
{
  "OPTION": VALUE,
  "OPTION": VALUE,
  "OPTION": VALUE
}
```

2\. All **VALUE**s are one of:

- **string**: ``"anythin inside of quotes"``
- **list**: ``["like", "a", "string", "broken", "into", "parts"]``
- **number**: ``5`` or ``3.1415``
- **boolean**: ``true`` for Yes or ``false`` for No

3\. Every **VALUE** must have a trailing comma (``,``), except the last must not.

Test-Driven Defaults
--------------------

Defaults were selected based on tests using [lowest-recommended hardware](hardware.md).

- [Encoder Documentation](https://trac.ffmpeg.org/wiki/Encode/H.264) (Upstream Defaults)

| Value                 | Purpose                                                            |
| --------------------- | ------------------------------------------------------------------ |
| ``-t <time>``         | Must be set on **each** input device                               |
| ``-tune zerolatency`` | This is **required** for software encoding on low-end cpu          |
| ``-bufsize 64M``      | Very large buffer to helps avoid processing spikes                 |
| ``-crf 23``           | Default is best; 25 reduces size by 30%, but 24 creates movement   |
| ``-maxrate 3M``       | Hard quality limit, based on **NO hardware encoding**              |
|                       | * 1080p recordings will see quality improvement up to about ``7M`` |
|                       | * Larger value creates larger files and XBUF for CPU-only encoding |
| ``-preset fast``      | Yield fewest XRUN errors                                           |
| ``-framerate 15``     | Fast enough to catch most movement                                 |

**Baseline Command:**

```
ffmpeg -y -loglevel warning -nostdin -nostats \
    -t 00:10:00 -f alsa -i plughw \
    -t 00:10:00 -f v4l2 -i /dev/video0 \
    -map 0:a -c:a pcm_s16le -ar 48000 -ac 1 -f wav - \
    -filter_complex [1:v]...[dtstamp] -map 0:a -map [dtstamp] \
    -c:a pcm_s16le -ar 48000 -ac 1 -c:v libx264 baseline.mkv \
    >/dev/null
```
