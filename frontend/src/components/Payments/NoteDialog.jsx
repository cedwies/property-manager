import React, { useState, useEffect } from 'react';
import './NoteDialog.css';

const NoteDialog = ({ note, onSave, onCancel }) => {
  const [noteText, setNoteText] = useState(note || '');

  useEffect(() => {
    setNoteText(note || '');
  }, [note]);

  const handleSave = () => {
    onSave(noteText);
  };

  return (
    <div className="note-dialog-overlay">
      <div className="note-dialog">
        <div className="note-dialog-header">
          <h3>Payment Note</h3>
          <button 
            className="close-button"
            onClick={onCancel}
            title="Close without saving"
          >
            ×
          </button>
        </div>
        
        <div className="note-dialog-content">
          <textarea
            className="note-textarea"
            value={noteText}
            onChange={(e) => setNoteText(e.target.value)}
            placeholder="Enter any relevant notes about this payment..."
            autoFocus
          />
        </div>
        
        <div className="note-dialog-footer">
          <button 
            className="button secondary"
            onClick={onCancel}
          >
            Cancel
          </button>
          <button 
            className="button primary"
            onClick={handleSave}
          >
            Save Note
          </button>
        </div>
      </div>
    </div>
  );
};

export default NoteDialog;