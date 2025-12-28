import React, { useState, useEffect } from 'react';
import api from '../services/api';

const CreateNoteModal = ({ onClose, onNoteCreated }) => {
    const [title, setTitle] = useState('');
    const [content, setContent] = useState('');
    const [loading, setLoading] = useState(false);

    const [availableTags, setAvailableTags] = useState([]);
    const [selectedTagIds, setSelectedTagIds] = useState([]);

    useEffect(() => {
        fetchTags();
    }, []);

    const fetchTags = async () => {
        try {
            const res = await api.get('/tags?tag_id=0&page_size=100');
            setAvailableTags(res.data || []);
        } catch (err) {
            console.error("Failed to fetch tags", err);
        }
    };

    const toggleTag = (tagId) => {
        setSelectedTagIds(prev =>
            prev.includes(tagId)
                ? prev.filter(id => id !== tagId)
                : [...prev, tagId]
        );
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        try {
            // 1. Create Note
            const noteRes = await api.post('/notes', { title, content });
            const newNote = noteRes.data;

            // 2. Link Tags
            if (newNote && newNote.note_id && selectedTagIds.length > 0) {
                // Execute in parallel
                await Promise.all(selectedTagIds.map(tagId =>
                    api.post('/note_tags', {
                        note_id: newNote.note_id,
                        tag_id: tagId
                    })
                ));
            }

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
            backdropFilter: 'blur(4px)', zIndex: 1000
        }}>
            <div style={{ background: 'white', padding: '2rem', borderRadius: '8px', width: '90%', maxWidth: '500px', maxHeight: '90vh', overflowY: 'auto' }}>
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

                    {/* Tag Selection */}
                    {availableTags.length > 0 && (
                        <div>
                            <label style={{ display: 'block', marginBottom: '0.5rem', fontWeight: 'bold', fontSize: '0.9rem' }}>Add Tags:</label>
                            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.5rem' }}>
                                {availableTags.map(tag => {
                                    const isSelected = selectedTagIds.includes(tag.tag_id);
                                    return (
                                        <div
                                            key={tag.tag_id}
                                            onClick={() => toggleTag(tag.tag_id)}
                                            style={{
                                                padding: '0.25rem 0.75rem',
                                                borderRadius: '999px',
                                                fontSize: '0.8rem',
                                                cursor: 'pointer',
                                                border: '1px solid',
                                                borderColor: isSelected ? 'var(--primary-color)' : '#ccc',
                                                background: isSelected ? 'var(--primary-color)' : 'transparent',
                                                color: isSelected ? 'white' : '#666'
                                            }}
                                        >
                                            {tag.name}
                                        </div>
                                    );
                                })}
                            </div>
                        </div>
                    )}

                    <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '1rem', marginTop: '1rem' }}>
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
