import json
import requests

server_key = "AAAAlUvY-Qk:APA91bHSeSeCE-sESkLYCa2sotvKXr2CPkPgLhVHwo07JKlgwFR0hbxMxM6or7TBUiNCdfwAxUhIcrlvu6d7MC76Flnj5XCxvu7J5VDZQ-y3tI1V_R29Q-4Qsfu_Ru5Pzp7oA9Np79Us"
mike_android_fcm = "etvClaI2QeiE-w9hWa75CK\:APA91bFgcvs5ngH68ve84wmCCQzDBdOrj69Ac1qrlLIVOS7fSEBCIBuaYX1Vr49ZN6ZIckg9xhFpy4xsasrehUUk16RS3rJlfvNDqPbgBCmHooVe8l8W8JGl9HLba9J1xcOELxOfCgqF"
url = "https://fcm.googleapis.com/fcm/send"

headers = {
    "Authorization": "key=" + server_key,
    "Content-Type": "application/json"
}

payload = {
    "to": mike_android_fcm,
    "data": {
        "notifee": json.dumps({
            "body": 'This message was sent via FCM!',
            "android": {
                "channelId": 'ring-sound',
                "actions": [{
                    "title": 'Mark as Read',
                    "pressAction": {
                        "id": 'read'
                    }
                }]
            }
        })
    }
}

r = requests.post(url, headers=headers, data=json.dumps(payload))
print(r.text)
