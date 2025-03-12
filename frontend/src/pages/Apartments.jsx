import React, { useState, useEffect } from 'react';
import ApartmentForm from '../components/Apartments/ApartmentForm';
import ApartmentsList from '../components/Apartments/ApartmentsList';
import ConfirmDialog from '../components/common/ConfirmDialog';
import './Apartments.css';

// Import Go backend functions
import { 
  GetAllHouses,
  GetAllApartments, 
  GetApartmentsByHouseID,
  CreateApartment, 
  UpdateApartment, 
  DeleteApartment 
} from '../../wailsjs/go/main/App';

const Apartments = () => {
  const [apartments, setApartments] = useState([]);
  const [houses, setHouses] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isFormVisible, setIsFormVisible] = useState(false);
  const [currentApartment, setCurrentApartment] = useState(null);
  const [error, setError] = useState(null);
  const [infoMessage, setInfoMessage] = useState(null);
  const [confirmDialog, setConfirmDialog] = useState({
    isOpen: false,
    message: '',
    apartmentId: null
  });
  
  // Fetch all apartments and houses on component mount
  useEffect(() => {
    loadData();
  }, []);
  
  const loadData = async () => {
    try {
      setIsLoading(true);
      setError(null);
      
      // Load houses first to be able to display house info for apartments
      const allHouses = await GetAllHouses();
      setHouses(allHouses || []);
      
      if (allHouses && allHouses.length > 0) {
        const allApartments = await GetAllApartments();
        setApartments(allApartments || []);
      } else {
        setApartments([]);
        setInfoMessage("Please add houses before adding apartments");
      }
    } catch (err) {
      console.error('Error loading data:', err);
      setError('Failed to load data. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };
  
  const handleAddApartment = () => {
    if (houses.length === 0) {
      setError('Please add at least one house before adding apartments');
      return;
    }
    
    setCurrentApartment(null);
    setIsFormVisible(true);
  };
  
  const handleEditApartment = (apartment) => {
    setCurrentApartment(apartment);
    setIsFormVisible(true);
  };
  
  const handleDeleteApartment = (id) => {
    setConfirmDialog({
      isOpen: true,
      message: 'Are you sure you want to delete this apartment?',
      apartmentId: id
    });
  };
  
  const confirmDelete = async () => {
    try {
      setError(null);
      await DeleteApartment(confirmDialog.apartmentId);
      setConfirmDialog({ isOpen: false, message: '', apartmentId: null });
      await loadData();
    } catch (err) {
      console.error('Error deleting apartment:', err);
      setError('Failed to delete apartment. Please try again.');
      setConfirmDialog({ isOpen: false, message: '', apartmentId: null });
    }
  };
  
  const cancelDelete = () => {
    setConfirmDialog({ isOpen: false, message: '', apartmentId: null });
  };
  
  const handleSaveApartment = async (apartmentData) => {
    try {
      setError(null);
      
      if (apartmentData.id) {
        // Update existing apartment
        await UpdateApartment(
          apartmentData.id,
          apartmentData.name,
          parseInt(apartmentData.houseId),
          apartmentData.size
        );
      } else {
        // Create new apartment
        await CreateApartment(
          apartmentData.name,
          parseInt(apartmentData.houseId),
          apartmentData.size
        );
      }
      
      // Reload apartments and close form
      await loadData();
      setIsFormVisible(false);
    } catch (err) {
      console.error('Error saving apartment:', err);
      setError(`Failed to ${apartmentData.id ? 'update' : 'add'} apartment: ${err.message || err}`);
    }
  };
  
  const handleCancelForm = () => {
    setIsFormVisible(false);
    setCurrentApartment(null);
  };
  
  return (
    <div className="apartments-page">
      <div className="apartments-header">
        <h1>Manage Apartments</h1>
        {!isFormVisible && (
          <button className="add-button" onClick={handleAddApartment}>
            Add New Apartment
          </button>
        )}
      </div>
      
      {error && <div className="error-banner">{error}</div>}
      {infoMessage && <div className="info-banner">{infoMessage}</div>}
      
      {isLoading ? (
        <div className="loading">Loading apartments...</div>
      ) : (
        <div className="apartments-content">
          {isFormVisible ? (
            <ApartmentForm 
              apartment={currentApartment} 
              houses={houses}
              onSave={handleSaveApartment} 
              onCancel={handleCancelForm} 
            />
          ) : (
            <ApartmentsList 
              apartments={apartments} 
              houses={houses}
              onEdit={handleEditApartment} 
              onDelete={handleDeleteApartment} 
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

export default Apartments;