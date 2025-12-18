import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Home from './pages/Home';
import Capture from './pages/Capture';
import Result from './pages/Result';
import Collection from './pages/Collection';
import Share from './pages/Share';
import LearningReport from './pages/LearningReport';

function App() {
  return (
    <BrowserRouter>
      <div className="min-h-screen bg-cloud-white">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/capture" element={<Capture />} />
          <Route path="/result" element={<Result />} />
          <Route path="/collection" element={<Collection />} />
          <Route path="/share/:shareId" element={<Share />} />
          <Route path="/report" element={<LearningReport />} />
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App;
