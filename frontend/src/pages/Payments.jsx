import React, { useState, useEffect } from 'react';
import HouseSelector from '../components/Payments/HouseSelector';
import PaymentRecordsDialog from '../components/Payments/PaymentRecordsDialog';
import './Payments.css';

import { GetAllHouses } from '../../wailsjs/go/main/App';

const Payments = () => {
  const [houses, setHouses] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showPaymentDialog, setShowPaymentDialog] = useState(false);
  const [selectedHouse, setSelectedHouse] = useState(null);

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

  const handleHouseSelect = (house) => {
    setSelectedHouse(house);
    setShowPaymentDialog(true);
  };

  const handleClosePaymentDialog = () => {
    setShowPaymentDialog(false);
    setSelectedHouse(null);
  };

  return (
    <div className="payments-page">
      <div className="payments-header">
        <h1>Payment Management</h1>
      </div>
      
      {error && <div className="error-banner">{error}</div>}
      
      {isLoading ? (
        <div className="loading">Loading houses...</div>
      ) : (
        <div className="payments-content">
          <div className="payments-description">
            <p>
              Select a house to check and manage payment records for all current tenants.
              You can record monthly payments, add notes, and lock records that are complete.
            </p>
          </div>
          
          <HouseSelector houses={houses} onSelectHouse={handleHouseSelect} />
        </div>
      )}
      
      {showPaymentDialog && selectedHouse && (
        <PaymentRecordsDialog 
          house={selectedHouse} 
          onClose={handleClosePaymentDialog}
        />
      )}
    </div>
  );
};

export default Payments;