import React, { useState, useEffect } from 'react';
import './App.css';

function App() {
  const [longUrl, setLongUrl] = useState('');
  const [shortUrl, setShortUrl] = useState('');
  const [urls, setUrls] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchUrls();
  }, []);

  const fetchUrls = async () => {
    try {
      const response = await fetch('/api/urls');
      const data = await response.json();
      setUrls(data || []);
    } catch (err) {
      console.error('Error fetching URLs:', err);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setShortUrl('');

    try {
      const response = await fetch('/api/urls', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ url: longUrl }),
      });

      if (!response.ok) {
        const errData = await response.json();
        throw new Error(errData.error || 'Failed to shorten URL');
      }

      const data = await response.json();
      const newShortUrl = `${window.location.origin}/${data.short_code}`;
      setShortUrl(newShortUrl);
      fetchUrls(); // Refresh the list
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>Che URL Shortener</h1>
      </header>
      <main>
        <form onSubmit={handleSubmit}>
          <input
            type="text"
            value={longUrl}
            onChange={(e) => setLongUrl(e.target.value)}
            placeholder="Enter a long URL"
            required
          />
          <button type="submit">Shorten URL</button>
        </form>
        {shortUrl && (
          <div className="result">
            <p>
              Short URL: <a href={shortUrl} target="_blank" rel="noopener noreferrer">{shortUrl}</a>
            </p>
          </div>
        )}
        {error && <p className="error">{error}</p>}
        <div className="history">
          <h2>URL History</h2>
          <table>
            <thead>
              <tr>
                <th>Short URL</th>
                <th>Original URL</th>
                <th>Clicks</th>
              </tr>
            </thead>
            <tbody>
              {urls.map((url) => (
                <tr key={url.short_code}>
                  <td>
                    <a href={`${window.location.origin}/${url.short_code}`} target="_blank" rel="noopener noreferrer">
                      {`${window.location.origin}/${url.short_code}`}
                    </a>
                  </td>
                  <td>{url.long_url}</td>
                  <td>{url.usage_count}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </main>
    </div>
  );
}

export default App;
