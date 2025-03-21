import React from 'react'
import ReactDOM from 'react-dom/client'
import './index.css'
import ChatRoom from './ChatRoom'

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <ChatRoom jwtToken={"example-jwt"} />
  </React.StrictMode>
)