import React, { useEffect, useState } from 'react';
import api from '../services/api';
import useAuth from '../hooks/useAuth';
import CreateNoteModal from '../components/CreateNoteModal';
import NoteDetailModal from '../components/NoteDetailModal';
import Sidebar from '../components/Sidebar';
import { Tag as TagIcon, Search } from 'lucide-react';
import '../styles/main.css';

const Dashboard = () => {
    const { user, logout } = useAuth();
    const [notes, setNotes] = useState([]);
    const [loadingNotes, setLoadingNotes] = useState(true);

    // Search state
    const [searchTerm, setSearchTerm] = useState('');

    // Modals
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [selectedNoteId, setSelectedNoteId] = useState(null); // For Detail Modal

    const [selectedTagId, setSelectedTagId] = useState(null);

    useEffect(() => {
        // Debounce search could be added here if needed, but for now direct fetch
        const timer = setTimeout(() => {
            fetchNotes();
        }, 300);
        return () => clearTimeout(timer);
    }, [selectedTagId, searchTerm]);

    const fetchNotes = async () => {
        setLoadingNotes(true);
        try {
            let url = '/notes?cursor=0&page_size=20';

            if (searchTerm) {
                // Should we respect tag filter + search? Backend API ListNotes only does OR logic (search or list).
                // Search takes precedence in my backend implementation of ListNotes
                url = `/notes?cursor=0&page_size=20&search=${encodeURIComponent(searchTerm)}`;
            } else if (selectedTagId) {
                url = `/tags/${selectedTagId}/notes`;
            }

            const res = await api.get(url);
            setNotes(res.data || []);
        } catch (err) {
            console.error("Failed to fetch notes", err);
            if (err.response && err.response.status === 404) {
                setNotes([]);
            }
        } finally {
            setLoadingNotes(false);
        }
    };

    return (
        <div style={{ minHeight: '100vh', background: 'var(--bg-secondary)', display: 'flex' }}>
            <Sidebar selectedTagId={selectedTagId} onSelectTag={(id) => { setSelectedTagId(id); setSearchTerm(''); }} />

            <div style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
                <header style={{ background: 'white', padding: '1rem 2rem', boxShadow: 'var(--shadow-sm)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <h1 style={{ fontSize: '1.5rem', fontWeight: 'bold', color: 'var(--primary-color)' }}>Notes</h1>

                    {/* Search Bar */}
                    <div style={{ flex: 1, maxWidth: '400px', margin: '0 2rem', position: 'relative' }}>
                        <Search size={18} style={{ position: 'absolute', left: '10px', top: '50%', transform: 'translateY(-50%)', color: '#999' }} />
                        <input
                            type="text"
                            placeholder="Search notes..."
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                            style={{
                                width: '100%', padding: '0.6rem 1rem 0.6rem 2.5rem',
                                borderRadius: '999px', border: '1px solid var(--border-color)',
                                outline: 'none', transition: 'border-color 0.2s'
                            }}
                            onFocus={(e) => e.target.style.borderColor = 'var(--primary-color)'}
                            onBlur={(e) => e.target.style.borderColor = 'var(--border-color)'}
                        />
                    </div>

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
                        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: '1.5rem' }}>
                            {notes.map(note => (
                                <div
                                    key={note.note_id}
                                    onClick={() => setSelectedNoteId(note.note_id)}
                                    style={{
                                        background: 'white', padding: '1.5rem', borderRadius: 'var(--radius-lg)',
                                        boxShadow: 'var(--shadow-sm)', border: '1px solid var(--border-color)',
                                        cursor: 'pointer', transition: 'transform 0.2s', position: 'relative',
                                        height: '200px', display: 'flex', flexDirection: 'column'
                                    }}
                                    onMouseEnter={e => e.currentTarget.style.transform = 'translateY(-2px)'}
                                    onMouseLeave={e => e.currentTarget.style.transform = 'translateY(0)'}
                                >
                                    <h3 style={{ marginBottom: '0.5rem', fontSize: '1.1rem', fontWeight: 'bold', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                                        {note.title?.String || note.title}
                                    </h3>

                                    {/* Tags Row */}
                                    <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.25rem', marginBottom: '0.5rem' }}>
                                        {note.tags && note.tags.slice(0, 3).map(tag => (
                                            <span key={tag.tag_id} style={{
                                                fontSize: '0.7rem', background: 'var(--bg-secondary)', color: 'var(--primary-color)',
                                                padding: '2px 6px', borderRadius: '4px', display: 'flex', alignItems: 'center', gap: '2px'
                                            }}>
                                                <TagIcon size={10} /> {tag.name}
                                            </span>
                                        ))}
                                        {note.tags && note.tags.length > 3 && <span style={{ fontSize: '0.7rem', color: '#999' }}>+{note.tags.length - 3}</span>}
                                    </div>

                                    {/* Content Preview (Truncated) */}
                                    <p style={{
                                        color: 'var(--text-muted)', fontSize: '0.9rem',
                                        flex: 1, overflow: 'hidden', textOverflow: 'ellipsis',
                                        display: '-webkit-box', WebkitLineClamp: 3, WebkitBoxOrient: 'vertical',
                                        wordBreak: 'break-word'
                                    }}>
                                        {note.content?.String || note.content}
                                    </p>

                                    <div style={{ marginTop: 'auto', fontSize: '0.75rem', color: '#ccc', paddingTop: '0.5rem' }}>
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
                        cursor: 'pointer',
                        border: 'none'
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

                {selectedNoteId && (
                    <NoteDetailModal
                        noteId={selectedNoteId}
                        onClose={() => setSelectedNoteId(null)}
                        onNoteUpdated={fetchNotes}
                    />
                )}
            </div>
        </div>
    )
}

export default Dashboard;
