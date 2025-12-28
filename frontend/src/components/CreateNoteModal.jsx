import React, { useState } from 'react';
import api from '../services/api';

const CreateNoteModal = ({ onClose, onNoteCreated }) => {
    const [title, setTitle] = useState('');
    const [content, setContent] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        try {
            await api.post('/notes', { title, content });
            onNoteCreated();
            onClose();
        } catch (err) {
            console.error("Failed to create note", err);
            alert("Failed to create note");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div style={{
            position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
            background: 'rgba(0,0,0,0.5)', display: 'flex', justifyContent: 'center', alignItems: 'center',
            backdropFilter: 'blur(4px)'
        }}>
            <div style={{ background: 'white', padding: '2rem', borderRadius: '8px', width: '90%', maxWidth: '500px' }}>
                <h2 style={{ marginBottom: '1rem' }}>Create New Note</h2>
                <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                    <input
                        type="text"
                        placeholder="Title"
                        value={title}
                        onChange={e => setTitle(e.target.value)}
                        required
                        style={{ padding: '0.5rem', borderRadius: '4px', border: '1px solid #ccc' }}
                    />
                    <textarea
                        placeholder="Content"
                        value={content}
                        onChange={e => setContent(e.target.value)}
                        required
                        rows={5}
                        style={{ padding: '0.5rem', borderRadius: '4px', border: '1px solid #ccc', fontFamily: 'inherit' }}
                    />
                    <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '1rem' }}>
                        <button type="button" onClick={onClose} style={{ color: '#666' }}>Cancel</button>
                        <button
                            type="submit"
                            disabled={loading}
                            style={{ background: '#4f46e5', color: 'white', padding: '0.5rem 1rem', borderRadius: '4px' }}
                        >
                            {loading ? 'Creating...' : 'Create'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default CreateNoteModal;
