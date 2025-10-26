Disturbance Tracker
===================

``DTrack`` is Surveilance Software that uses Machine Learning (AI/ML) to
learn about and report on local disturbances, using low-end hardware.

**Documentation:** https://mtecknology.github.io/dtrack/

Background
----------

Various excuses were being used to ignore reports of nuisance barking. The only
solution was to wait for at least one hour of nuisance barking across two days
and then build a detailed log that tracked every minute that a dog barked.

Similar solutions to this problem existed, but generic solutions fell horribly
short (confusing wind with dog barks) and targeted solutions required high-cost
resources (AWS + always-on internet) to operate.

``Disturbance Tracker (DTrack)`` solves this problem without external resources.

How It Works
------------

  1. Set up ``Monitoring Device`` (like a Raspberry Pi)
  2. Collect some initial recordings
  3. Train a model from collected recordings
  4. Use trained model for automatic detection

  .. image:: https://raw.githubusercontent.com/MTecknology/dtrack/refs/heads/master/.github/images/workflow.webp
     :alt: Disturbance Tracker Workflow

Monitoring Device
~~~~~~~~~~~~~~~~~

DTrack is designed to run on a ``Raspberry Pi v5``. It does not require any
internet connection or subscription-style service.


Configuration Options
---------------------

Configuration is stored in a ``JSON`` file, called ``config.json``. A sample can
be copied from ``example_config.json``.

  - List audio devices: ``ffmpeg -loglevel warning -sources alsa``
    + Look for: ``Hardware device with all software conversions``
  - List video devices: ``v4l2-ctl --list-devices --all``

Full configuration options are available using ``./dtrack --help``.

**Test-Driven Defaults:**

  - `Encoder Documentation <https://trac.ffmpeg.org/wiki/Encode/H.264>`__
  - ``-t <time>`` **must** be set prior to each input
  - ``-tune zerolatency`` is **required** for software encoding on low-end cpu
  - ``-bufsize 64M``: Very large buffer to helps avoid processing spikes
  - ``-crf 23``: Default is best; 25 reduces size by 30%, but 24 shows significant loss
  - ``-maxrate 3M``: Hard quality limit, based on **NO hardware encoding**

    + Benefits can be seen up to ``6M``
    + Larger value creates larger files and XBUF for CPU-only encoding

  - ``-preset fast``: Yield fewest XRUN errors
  - ``-framerate 15``: Fast enough to catch most movement

Baseline Command::

    ffmpeg -y -loglevel warning -nostdin -nostats \
        -t 00:10:00 -f alsa -i plughw \
        -t 00:10:00 -f v4l2 -i /dev/video0 \
        -map 0:a -c:a pcm_s16le -ar 48000 -ac 1 -f wav - \
        -filter_complex [1:v]...[dtstamp] -map 0:a -map [dtstamp] \
        -c:a pcm_s16le -ar 48000 -ac 1 -c:v libx264 -tune zerolatency baseline.mkv \
        >/dev/null
