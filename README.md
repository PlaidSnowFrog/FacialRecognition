
# Facial Recognition

This is a Facial Recognition project that logs at what time faces were detected

## Installation

```
$ curl -L -O https://github.com/PlaidSnowFrog/FacialRecognition/releases/download/v1.0.0/FacialRecognition
&& mkdir -p data/
&& curl -L https://raw.githubusercontent.com/opencv/opencv/master/data/haarcascades/haarcascade_frontalface_default.xml -o data/haarcascade_frontalface_default.xml
&& curl -L https://raw.githubusercontent.com/opencv/opencv/master/data/haarcascades/haarcascade_profileface.xml -o data/haarcascade_profileface.xml
&& curl -L https://raw.githubusercontent.com/opencv/opencv/master/data/haarcascades/haarcascade_eye.xml -o data/haarcascade_eye.xml
&& touch log.txt
```
