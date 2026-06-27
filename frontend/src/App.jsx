import React, { useState } from 'react';
import './App.css';

function App() {
  const [file, setFile] = useState(null);
  const [analyzing, setAnalyzing] = useState(false);
  const [result, setResult] = useState(null);

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleUpload = async (e) => {
    e.preventDefault();
    if (!file) return;

    setAnalyzing(true);
    const formData = new FormData();
    formData.append('file', file);

    try {
      const response = await fetch('/api/analyze', {
        method: 'POST',
        body: formData,
      });
      const data = await response.json();
      setResult(data);
    } catch (error) {
      console.error('Upload failed:', error);
      setResult({ error: 'Upload failed' });
    } finally {
      setAnalyzing(false);
    }
  };

  const handleReport = async () => {
    if (!result || !result.patternId) return;

    try {
      await fetch('/api/report', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ patternId: result.patternId }),
      });
      alert('記録しました');
      setResult(null);
      setFile(null);
    } catch (error) {
      console.error('Report failed:', error);
    }
  };

  return (
    <div className="container">
      <h1>闇バイト判定アプリ</h1>
      <p>スクショを投げるだけで、同じ文言パターンが何件記録されているかを返します。</p>

      {!result ? (
        <form onSubmit={handleUpload}>
          <input
            type="file"
            accept="image/*"
            onChange={handleFileChange}
            disabled={analyzing}
          />
          <button type="submit" disabled={!file || analyzing}>
            {analyzing ? '分析中...' : '分析'}
          </button>
        </form>
      ) : (
        <div className="result">
          <h2>分析結果</h2>
          {result.error ? (
            <p style={{ color: 'red' }}>{result.error}</p>
          ) : (
            <>
              <p>
                <strong>このパターンは {result.count} 件記録されています</strong>
              </p>
              <p>{result.normalizedText}</p>
              <button onClick={handleReport}>記録する</button>
              <button onClick={() => setResult(null)}>キャンセル</button>
            </>
          )}
        </div>
      )}
    </div>
  );
}

export default App;
