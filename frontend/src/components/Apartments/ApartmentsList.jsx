import React, { useState, useEffect } from 'react';
import './ApartmentsList.css';

const ApartmentsList = ({ apartments, houses, onEdit, onDelete }) => {
  const [filteredApartments, setFilteredApartments] = useState([]);
  const [selectedHouseId, setSelectedHouseId] = useState('all');

  useEffect(() => {
    if (apartments && apartments.length > 0) {
      filterApartments(selectedHouseId);
    } else {
      setFilteredApartments([]);
    }
  }, [apartments, selectedHouseId]);

  const filterApartments = (houseId) => {
    if (houseId === 'all') {
      setFilteredApartments(apartments);
    } else {
      const filtered = apartments.filter(
        apartment => apartment.houseId.toString() === houseId
      );
      setFilteredApartments(filtered);
    }
  };

  const handleFilterChange = (e) => {
    setSelectedHouseId(e.target.value);
  };

  const formatSize = (size) => {
    return typeof size === 'number' ? size.toFixed(2).replace(/\.00$/, '') : '';
  };

  if (!apartments || apartments.length === 0) {
    return (
      <div className="apartments-list empty">
        <p>No apartments found. Add your first apartment using the form.</p>
      </div>
    );
  }

  return (
    <div className="apartments-list">
      <div className="apartments-header">
        <h2>Your Apartments</h2>
        <div className="filter-container">
          <label htmlFor="house-filter">Filter by House:</label>
          <select
            id="house-filter"
            value={selectedHouseId}
            onChange={handleFilterChange}
          >
            <option value="all">All Houses</option>
            {houses && houses.map(house => (
              <option key={house.id} value={house.id.toString()}>
                {house.name}
              </option>
            ))}
          </select>
        </div>
      </div>

      <div className="list-container">
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>House</th>
              <th style={{ textAlign: 'right' }}>Size (m²)</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {filteredApartments.map((apartment) => (
              <tr key={apartment.id}>
                <td>{apartment.name}</td>
                <td>
                  {apartment.house ? apartment.house.name : 
                   (houses.find(h => h.id === apartment.houseId)?.name || 'Unknown')}
                </td>
                <td className="apartment-size">{formatSize(apartment.size)} m²</td>
                <td className="actions">
                  <button 
                    className="button edit" 
                    onClick={() => onEdit(apartment)}
                    title="Edit apartment"
                  >
                    Edit
                  </button>
                  <button 
                    className="button delete" 
                    onClick={() => onDelete(apartment.id)}
                    title="Delete apartment"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default ApartmentsList;