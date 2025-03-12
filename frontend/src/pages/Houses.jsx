import React, { useState, useEffect } from 'react';
import HouseForm from '../components/Houses/HouseForm';
import HousesList from '../components/Houses/HousesList';
import ConfirmDialog from '../components/common/ConfirmDialog';
import './Houses.css';

// Import Go backend functions
import { CreateHouse, GetAllHouses, UpdateHouse, DeleteHouse } from '../../wailsjs/go/main/App';

const Houses = () => {
  const [houses, setHouses] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isFormVisible, setIsFormVisible] = useState(false);
  const [currentHouse, setCurrentHouse] = useState(null);
  const [error, setError] = useState(null);
  const [confirmDialog, setConfirmDialog] = useState({
    isOpen: false,
    message: '',
    houseId: null
  });
  
  // Fetch all houses on component mount
  useEffect(() => {
    loadHouses();
  }, []);
  
  const loadHouses = async () => {
    try {
      setIsLoading(true);
      setError(null);
      
      const allHouses = await GetAllHouses();
      setHouses(allHouses || []);
    } catch (err) {
      console.error('Error loading houses:', err);
      setError('Failed to load houses. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };
  
  const handleAddHouse = () => {
    setCurrentHouse(null);
    setIsFormVisible(true);
  };
  
  const handleEditHouse = (house) => {
    setCurrentHouse(house);
    setIsFormVisible(true);
  };
  
  const handleDeleteHouse = (id) => {
    setConfirmDialog({
      isOpen: true,
      message: 'Are you sure you want to delete this house?',
      houseId: id
    });
  };
  
  const confirmDelete = async () => {
    try {
      setError(null);
      await DeleteHouse(confirmDialog.houseId);
      setConfirmDialog({ isOpen: false, message: '', houseId: null });
      await loadHouses();
    } catch (err) {
      console.error('Error deleting house:', err);
      setError('Failed to delete house. Please try again.');
      setConfirmDialog({ isOpen: false, message: '', houseId: null });
    }
  };
  
  const cancelDelete = () => {
    setConfirmDialog({ isOpen: false, message: '', houseId: null });
  };
  
  const handleSaveHouse = async (houseData) => {
    try {
      setError(null);
      
      if (houseData.id) {
        // Update existing house
        await UpdateHouse(
          houseData.id,
          houseData.name,
          houseData.street,
          houseData.number,
          houseData.country,
          houseData.zipCode,
          houseData.city
        );
      } else {
        // Create new house
        await CreateHouse(
          houseData.name,
          houseData.street,
          houseData.number,
          houseData.country,
          houseData.zipCode,
          houseData.city
        );
      }
      
      // Reload houses and close form
      await loadHouses();
      setIsFormVisible(false);
    } catch (err) {
      console.error('Error saving house:', err);
      setError(`Failed to ${houseData.id ? 'update' : 'add'} house: ${err.message || err}`);
    }
  };
  
  const handleCancelForm = () => {
    setIsFormVisible(false);
    setCurrentHouse(null);
  };
  
  return (
    <div className="houses-page">
      <div className="houses-header">
        <h1>Manage Houses</h1>
        {!isFormVisible && (
          <button className="add-button" onClick={handleAddHouse}>
            Add New House
          </button>
        )}
      </div>
      
      {error && <div className="error-banner">{error}</div>}
      
      {isLoading ? (
        <div className="loading">Loading houses...</div>
      ) : (
        <div className="houses-content">
          {isFormVisible ? (
            <HouseForm 
              house={currentHouse} 
              onSave={handleSaveHouse} 
              onCancel={handleCancelForm} 
            />
          ) : (
            <HousesList 
              houses={houses} 
              onEdit={handleEditHouse} 
              onDelete={handleDeleteHouse} 
            />
          )}
        </div>
      )}
      
      <ConfirmDialog 
        isOpen={confirmDialog.isOpen}
        message={confirmDialog.message}
        onConfirm={confirmDelete}
        onCancel={cancelDelete}
      />
    </div>
  );
};

export default Houses;