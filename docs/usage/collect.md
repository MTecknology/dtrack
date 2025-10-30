Data Collection
===============

Automatic detection cannot happen without a trained model, and a model cannot
be trained without some initial data, recorded from our detection hardware.

This loop is primed using audio data taken from **your** recording hardware.

!!! warning "Critical!"
    Changing the recording device, capture band, or even band can have a dramatic
    result on the ability of a trained model to continue accurate detection.


1\. Begin collecting recordings:
Begin continuous recording with:
```sh
    dtrack -a monitor -v
```

Stop recording with:
```sh
    # From the same terminal session to cancel when current recording finishes
    Ctrl+C

    # Press a second time to disregard current recording and exit immediately
    Ctrl+C
```

Recordings will be saved to ``./_workspace/rotating/``.

!!! note "Demonstration Note:"

    Clapping hands together is a great demonstration exercise. This can be set
    in `config.yml` with `inspect_models: [clap]`.
