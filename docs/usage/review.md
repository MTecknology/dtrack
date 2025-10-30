Data Analysis
=============

During initial data collection, it can be useful to set `record_duration:` to
2-5 minutes and then rename each recording as they complete, using the following
as an example:

```text
    2024-08-10_13:54:00_TRAIN-clap.mkv
    2024-08-10_13:58:17_TRAIN-clap.mkv
    2024-08-10_14:00:26_TEST-nomatch.mkv
    2024-08-10_14:02:38_TEST-nomatch.mkv
    2024-08-10_14:04:57_TEST-nomatch.mkv
    2024-08-10_14:07:20_TEST-nomatch.mkv
    2024-08-10_14:09:45_TEST-clap.mkv
    2024-08-10_14:17:07_TEST-nomatch.mkv
    2024-08-10_14:19:35_TEST-nomatch.mkv
    2024-08-10_14:22:02_TEST-clap.mkv
```

- ``TRAIN`` files were created with many variations of the sound being searched
    for, using different background noise, volumes, etc. These files will be cut
    into 1-second clips for model training.
- ``TEST`` files include variation, but may have only one instance of a sound
    being searched for--one needle in the haystack. These will be used to test the
    quality of each ML iteration.

Once data is collected, it can be retrieved from `_workspace/rotating` on the
[recording device](collect.md) and copied to the same `_workspace/rotating`
location on the [device used for training](train.md).

**Testing Data:**

Test data is essentially the same as training data, except it is collected with
the intent of being used only for testing.

Follow the process for tagging and then move data to
`_workspace/test/<model>/` or `_workspace/test/nomatch/`:

Ultimately, these videos will be used to determine the accuracy of each model.

**Training Data:**

In order to determine if something `is` or `is not`, the source audio must
be broken up into short consumable segments and segments matching the target
model must be reviewed and saved (tagged) manually.

!!! note "Project Timing"
    - DTrack is designed for generating reports.
    - Report granularity uses 1-minute cycles.
      + 1 clap or 999 claps within 1 minute is logged as one hit.
    - Each recording is broken into 2-second clips.
    - Each clip overlaps the next by 1 second to prevent dead zones

Open and review captured (from `rotating/`) using the inspection tool:
```sh
    dtrack -a review
```

The `review` option provides a GUI to help simplify the process of reviewing
and tagging 1-second clips.

Keyboard Shortcuts:

  - Left/Right: Navigate 1 frame left or right
  - PgUp/PgDn: Navigate 60 frames left or right
  - Home/End: Navigate to start or end
  - Up: Replay audio clip
