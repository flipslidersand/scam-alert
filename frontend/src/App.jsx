import React, { useState } from 'react';
import { analyzeImage, reportPattern } from './api/client';
import './App.css';

function App() {
  const [file, setFile] = useState(null);
  const [analyzing, setAnalyzing] = useState(false);
  const [result, setResult] = useState(null);
  const [error, setError] = useState(null);

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
    setError(null);
  };

  const handleUpload = async (e) => {
    e.preventDefault();
    if (!file) return;

    setAnalyzing(true);
    setError(null);

    try {
      const data = await analyzeImage(file);
      if (data.error) {
        setError(data.error);
        setResult(null);
      } else {
        setResult(data);
      }
    } catch (err) {
      setError(err.error || 'Failed to analyze image');
      setResult(null);
    } finally {
      setAnalyzing(false);
    }
  };

  const handleReport = async () => {
    if (!result || !result.patternId) return;

    try {
      await reportPattern(result.patternId);
      setResult(null);
      setFile(null);
      setError(null);
      // 成功メッセージはシンプルに次の分析に進む
    } catch (err) {
      setError(err.error || 'Failed to report pattern');
    }
  };

  const handleCancel = () => {
    setResult(null);
    setFile(null);
    setError(null);
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
          {error && <p className="error">{error}</p>}
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
              <p className="normalized-text">{result.normalizedText}</p>
              <div className="button-group">
                <button className="btn-primary" onClick={handleReport}>
                  記録する
                </button>
                <button className="btn-secondary" onClick={handleCancel}>
                  キャンセル
                </button>
              </div>
            </>
          )}
        </div>
      )}
    </div>
  );
}

export default App;
