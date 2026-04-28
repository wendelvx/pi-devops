import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import SetupScreen from './screens/SetupScreen';
import TelaoScreen from './screens/TelaoScreen';

export default function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<SetupScreen />} />
        {/* A rota do Telão recebe o código da sala na URL */}
        <Route path="/arena/:roomCode" element={<TelaoScreen />} />
      </Routes>
    </Router>
  );
}