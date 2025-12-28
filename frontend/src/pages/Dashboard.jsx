import React, { useEffect, useState } from 'react';
import api from '../services/api';
import useAuth from '../hooks/useAuth';
import CreateNoteModal from '../components/CreateNoteModal';
import Sidebar from '../components/Sidebar';
import '../styles/main.css';

const Dashboard = () => {
    const { user, logout } = useAuth();
    const [notes, setNotes] = useState([]);
    const [loadingNotes, setLoadingNotes] = useState(true);
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [selectedTagId, setSelectedTagId] = useState(null);

    useEffect(() => {
        fetchNotes();
    }, [selectedTagId]);

    const fetchNotes = async () => {
        setLoadingNotes(true);
        try {
            let url = '/notes?cursor=0&page_size=20';
            if (selectedTagId) {
                url = `/tags/${selectedTagId}/notes`;
            }

            const res = await api.get(url);
            setNotes(res.data || []);
        } catch (err) {
            console.error("Failed to fetch notes", err);
            // If 404 on tag notes (no notes), set empty
            if (err.response && err.response.status === 404) {
                setNotes([]);
            }
        } finally {
            setLoadingNotes(false);
        }
    };

    return (
        <div style={{ minHeight: '100vh', background: 'var(--bg-secondary)', display: 'flex' }}>
            <Sidebar selectedTagId={selectedTagId} onSelectTag={setSelectedTagId} />

            <div style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
                <header style={{ background: 'white', padding: '1rem 2rem', boxShadow: 'var(--shadow-sm)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <h1 style={{ fontSize: '1.5rem', fontWeight: 'bold', color: 'var(--primary-color)' }}>Notes</h1>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                        <span>Welcome, <b>{user?.username}</b></span>
                        <button onClick={logout} style={{ padding: '0.5rem 1rem', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-md)', background: 'white' }}>
                            Logout
                        </button>
                    </div>
                </header>

                <main style={{ padding: '2rem', flex: 1, overflowY: 'auto' }}>
                    {loadingNotes ? (
                        <p>Loading notes...</p>
                    ) : notes.length === 0 ? (
                        <div style={{ textAlign: 'center', marginTop: '3rem', color: 'var(--text-muted)' }}>
                            <p>No notes found.</p>
                            <button onClick={() => setShowCreateModal(true)} style={{ marginTop: '1rem', color: 'var(--primary-color)', fontWeight: 'bold' }}>Create one?</button>
                        </div>
                    ) : (
                        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))', gap: '1.5rem' }}>
                            {notes.map(note => (
                                <div key={note.note_id} style={{ background: 'white', padding: '1.5rem', borderRadius: 'var(--radius-lg)', boxShadow: 'var(--shadow-sm)', border: '1px solid var(--border-color)' }}>
                                    <h3 style={{ marginBottom: '0.5rem', fontSize: '1.1rem' }}>{note.title?.String || note.title}</h3>
                                    <p style={{ color: 'var(--text-muted)', whiteSpace: 'pre-wrap' }}>{note.content?.String || note.content}</p>
                                    <div style={{ marginTop: '1rem', fontSize: '0.875rem', color: '#9ca3af' }}>
                                        {note.created_at ? new Date(note.created_at.Time || note.created_at).toLocaleDateString() : ''}
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}
                </main>

                <button
                    onClick={() => setShowCreateModal(true)}
                    style={{
                        position: 'fixed',
                        bottom: '2rem',
                        right: '2rem',
                        background: 'var(--primary-color)',
                        color: 'white',
                        width: '3.5rem',
                        height: '3.5rem',
                        borderRadius: '50%',
                        fontSize: '2rem',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        boxShadow: 'var(--shadow-lg)',
                        cursor: 'pointer'
                    }}
                >
                    +
                </button>

                {showCreateModal && (
                    <CreateNoteModal
                        onClose={() => setShowCreateModal(false)}
                        onNoteCreated={fetchNotes}
                    />
                )}
            </div>
        </div>
    )
}

export default Dashboard;
