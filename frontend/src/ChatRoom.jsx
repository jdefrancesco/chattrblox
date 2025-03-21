import React, { useEffect, useRef, useState } from 'react'

// Front end chat room code.
const ChatRoom = ({ jwtToken }) => {
  const wsRef = useRef(null)
  const pcRef = useRef(null)
  const localVideo = useRef(null)
  const remoteVideo = useRef(null)
  const [partnerID, setPartnerID] = useState(null)


  useEffect(() => {
    // Connect to websocket.
    wsRef.current = new WebSocket('ws://localhost:8080/ws', jwtToken)

    wsRef.current.onmessage = async (event) => {
      const msg = JSON.parse(event.data)

      // Handle message types
      switch (msg.type) {
        case 'match_found':
          setPartnerID(msg.data)
          await startWebRTC()
          break
        case 'signal_offer':
          await handleOffer(msg)
          break
        case 'signal_answer':
          await pcRef.current.setRemoteDescription(new RTCSessionDescription(JSON.parse(msg.data)))
          break
        case 'ice_candidate':
          await pcRef.current.addIceCandidate(new RTCIceCandidate(JSON.parse(msg.data)))
          break
        case 'disconnect':
          endCall()
          break
      }
    }

    return () => wsRef.current?.close()
  }, [])

  // Start WebRTC. Use STUN to setup p2p connection.
  const startWebRTC = async () => {
    pcRef.current = new RTCPeerConnection({
      iceServers: [{ urls: 'stun:stun.l.google.com:19302' }],
    })

    pcRef.current.onicecandidate = (e) => {
      if (e.candidate && partnerID) {
        wsRef.current.send(JSON.stringify({
          type: 'ice_candidate',
          data: JSON.stringify(e.candidate),
          to: partnerID
        }))
      }
    }

    pcRef.current.ontrack = (e) => {
      remoteVideo.current.srcObject = e.streams[0]
    }

    const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true })
    stream.getTracks().forEach(track => pcRef.current.addTrack(track, stream))
    localVideo.current.srcObject = stream

    const offer = await pcRef.current.createOffer()
    await pcRef.current.setLocalDescription(offer)

    wsRef.current.send(JSON.stringify({
      type: 'signal_offer',
      data: JSON.stringify(offer),
      to: partnerID
    }))
  }

  // Handle offer.
  const handleOffer = async (msg) => {
    pcRef.current = new RTCPeerConnection({ iceServers: [{ urls: 'stun:stun.l.google.com:19302' }] })
    pcRef.current.onicecandidate = (e) => {
      if (e.candidate) {
        wsRef.current.send(JSON.stringify({
          type: 'ice_candidate',
          data: JSON.stringify(e.candidate),
          to: msg.to
        }))
      }
    }
    pcRef.current.ontrack = (e) => {
      remoteVideo.current.srcObject = e.streams[0]
    }

    const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true })
    stream.getTracks().forEach(track => pcRef.current.addTrack(track, stream))
    localVideo.current.srcObject = stream

    await pcRef.current.setRemoteDescription(new RTCSessionDescription(JSON.parse(msg.data)))
    const answer = await pcRef.current.createAnswer()
    await pcRef.current.setLocalDescription(answer)

    wsRef.current.send(JSON.stringify({
      type: 'signal_answer',
      data: JSON.stringify(answer),
      to: msg.to
    }))
    setPartnerID(msg.to)
  }


  // Next button handler
  const handleNext = () => {
    if (wsRef.current && partnerID) {
      wsRef.current.send(JSON.stringify({
        type: 'disconnect',
        to: partnerID,
      }))
    }

    endCall()
    // Runqueue for next match here!
    wsRef.current.send(JSON.stringify({ type: 'ready' }))
  }

  // Report button handler
  const handleReport = () => {
    if (wsRef.current && partnerID) {
      wsRef.current.send(JSON.stringify({
        type: 'report',
        to: partnerID,
        data: 'abuse',
      }))
      alert('User reported. Thank you.')
    }
  }

  // Reset state.
  const endCall = () => {
    pcRef.current?.getSenders().forEach(s => pcRef.current.removeTrack(s))
    pcRef.current?.close()
    pcRef.current = null
    setPartnerID(null)
    remoteVideo.current.srcObject = null
  }

  // Add chatroom element on page. Next and Report buttons right after it.
  return (
    <>
    <div className="chatroom-container">
        <div className="chatroom">
          <video ref={localVideo} autoPlay muted />
          <video ref={remoteVideo} autoPlay />
        </div>

        <div className="flex flex-col space-y-2 mt-4">
            <button onClick={handleNext} className="bg-yellow-500 hover:bg-yellow-600 p-2 rounded-xl">Next</button>
            <button onClick={handleReport} className="bg-red-500 hover:bg-red-600 p-2 rounded-xl">Report</button>
        </div>
    </div>

    <div className="h-screen w-screen bg-gray-900 text-white flex flex-col items-center justify-center p-4 space-y-4">
      <div className="grid grid-cols-2 gap-4 w-full max-w-5xl">
        <div className="flex flex-col items-center">
          <h2 className="text-lg mb-2">You</h2>
          <video ref={localVideo} className="rounded-xl shadow w-full" autoPlay muted />
        </div>
        <div className="flex flex-col items-center">
          <h2 className="text-lg mb-2">Stranger</h2>
          <video ref={remoteVideo} className="rounded-xl shadow w-full" autoPlay />
        </div>
      </div>

      <div className="w-full max-w-5xl grid grid-cols-3 gap-4">
        <div className="col-span-2">
          <div className="bg-gray-800 rounded-xl p-4 h-60 overflow-y-auto">
            {messages.map((msg, idx) => (
              <div key={idx} className={`mb-1 ${msg.from === 'me' ? 'text-blue-300' : 'text-pink-300'}`}>
                <span className="font-semibold">{msg.from === 'me' ? 'You' : 'Stranger'}:</span> {msg.text}
              </div>
            ))}
          </div>
          <input
            type="text"
            placeholder="Type a message..."
            className="mt-2 w-full p-2 rounded bg-gray-700 text-white outline-none"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && sendMessage()}
          />
        </div>

        <div className="flex flex-col space-y-2">
          <button onClick={handleNext} className="bg-yellow-500 hover:bg-yellow-600 p-2 rounded-xl">Next</button>
          <button onClick={handleReport} className="bg-red-500 hover:bg-red-600 p-2 rounded-xl">Report</button>
        </div>
      </div>
    </div>
    </>
  )
}

export default ChatRoom
