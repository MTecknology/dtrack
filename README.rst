Disturbance Tracker
===================

``DTrack`` is Surveilance Software that uses Machine Learning (AI/ML) to
learn about and report on local disturbances, without nvidia or aws.

**Documentation:** https://mtecknology.github.io/dtrack/

Background
----------

Various excuses were being used to ignore reports of nuisance barking. The only
solution was to wait for at least one hour of nuisance barking across two days
and then build a detailed log that tracked every minute that a dog barked.

Similar solutions to this problem existed, but generic solutions fell horribly
short (confusing wind with dog barks) and targeted solutions required high-cost
resources (AWS + always-on internet) to operate.

``DTrack`` attempts to solve the same problem without external resources.

How It Works
------------

  .. image:: https://raw.githubusercontent.com/MTecknology/dtrack/refs/heads/master/.github/images/workflow.webp
     :alt: Disturbance Tracker Workflow

1. Set up ``Monitoring Device``
2. Collect some initial recordings
3. Manually review recordings and tag nuisance noises
4. Train a model (i.e. "Teach AI")
5. Use trained model for automatic detection

Monitoring Device
~~~~~~~~~~~~~~~~~

DTrack is designed to run on a ``Raspberry Pi v5``. It does not require any
internet connection or subscription-style service.


Configuration Options
---------------------

Configuration is stored in a ``JSON`` file, called ``config.json``. A sample can
be copied from ``example_config.json``.

Full configuration options are available using ``./dtrack --help``.

**TODO:** ``docs -> github -> detailed config option docs``

**record_inspect_segment**: This is super challenging to automate.

  1. Create an initial recording

  2. Inspect recording:

     .. code-block:: text

        ffprobe -v error -select_streams a:0 \
            -show_entries stream=sample_rate,channels,bits_per_raw_sample
            -of default _workspace/recordings/<SINGLE-RECORDING>.mkv

        [STREAM]
        sample_rate=48000
        channels=2
        bits_per_raw_sample=N/A
        [/STREAM]

  3. If ``bits_per_raw_sample=NA``, then use ``16``
  4. record_inspect_segment: ``sample_rate * channels * (bits_per_raw_sample / 8)``
  5. Example: ``N/A -> 16``, ``record_inspect_segment = 48000 * 2 * (16/8) = 192000``
