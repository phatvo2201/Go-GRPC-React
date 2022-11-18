import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import { AuthProvider } from './context/AuthProvider';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import VerifyScreen from './components/VerifyScreen'

ReactDOM.render(
  <React.StrictMode>
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/*" element={<App />} />
        </Routes>
      </AuthProvider>
      <Routes>
        <Route path="/verifyemail/:verificationCode" element={<VerifyScreen />} />

      </Routes>
    </BrowserRouter>
  </React.StrictMode>,
  document.getElementById('root')
);