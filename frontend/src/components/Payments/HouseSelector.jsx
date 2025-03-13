import React from 'react';
import './HouseSelector.css';

const HouseSelector = ({ houses, onSelectHouse }) => {
  if (!houses || houses.length === 0) {
    return (
      <div className="house-selector-empty">
        <p>No houses found. Please add houses before managing payments.</p>
      </div>
    );
  }

  return (
    <div className="house-selector">
      <h2>Select a House to Check Payments</h2>
      <div className="house-grid">
        {houses.map((house) => (
          <div 
            key={house.id} 
            className="house-card" 
            onClick={() => onSelectHouse(house)}
          >
            <div className="house-card-name">{house.name}</div>
            <div className="house-card-address">
              {house.street} {house.number}, {house.zipCode} {house.city}
            </div>
            <div className="house-card-action">
              <button className="check-payments-button">
                Check Payments
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default HouseSelector;