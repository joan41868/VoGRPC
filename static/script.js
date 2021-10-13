let recorder;
let mediaStream;

const socket = io("ws://localhost:3000", {transports: ['websocket']});

socket.on('connect', () => {
    console.log("connected to backend.");

    socket.emit("subscribe", {
        roomID: "test",
        sender: "test"
    });
});

socket.on('receiveVoiceMessage',  async (vm) => {
    console.log('Received voice msg: ', vm);
    const raw = vm.data.split("base64,")[1];
    const blob = await b64toBlob(raw, vm.type);
    console.log(blob);
    const audioUrl = URL.createObjectURL(blob);
    const audio = new Audio(audioUrl);
    await audio.play();
});

async function onDataWS(blobEvt) {
    const text = await blobToBase64(blobEvt.data);
    const vm = {
        type: 'audio/webm;codecs=opus',
        data: text,
        sender: "test",
        roomID: "test"
    };
    console.log("Send: ", vm);
    socket.emit("sendVoiceMessage", vm);
}

async function startRec() {
    if (navigator.mediaDevices && navigator.mediaDevices.getUserMedia) {
        mediaStream = await navigator.mediaDevices.getUserMedia(
            // constraints - only audio needed for this app
            {
                audio: true
            });
        recorder = new MediaRecorder(mediaStream);
        recorder.start();
        recorder.ondataavailable = onDataWS;
    } else {
        console.log('getUserMedia not supported on your browser!');
    }
}

function stopRec() {
    recorder.stop();
    mediaStream.getTracks().forEach(t => {
        t.stop();
    });
}
