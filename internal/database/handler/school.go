package handler

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang-microservice-sekolah/internal/database"
	"golang-microservice-sekolah/internal/database/model"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func GetSchoolsHandler(w http.ResponseWriter, r *http.Request) {
	dbService := database.New()
	db := dbService.ToGormDB()
	
	var schools []model.School
	
	if err := db.Model(&model.School{}).
		Select("uuid, name, kode_provinsi, kode_kab_kota, kode_kecamatan, npsn, bentuk, status, alamat_jalan, lintang, bujur").
		Find(&schools).Error; err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(schools); err != nil {
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
	}
}

func GetSchoolByUuidHandler(w http.ResponseWriter, r *http.Request) {
	uuidParam := chi.URLParam(r, "uuid")
	
	dbService := database.New()
	db := dbService.ToGormDB()
	
	var school model.School
	
	if err := db.Model(&model.School{}).
		Select("uuid, name, kode_provinsi, kode_kab_kota, kode_kecamatan, npsn, bentuk, status, alamat_jalan, lintang, bujur").
		Where("uuid = ?", uuidParam).First(&school).Error; err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(school); err != nil {
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
	}
}

func CreateSchoolHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	
	var newSchool model.School
	
	if err := decoder.Decode(&newSchool); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	newSchool.UUID = uuid.New()
	
	dbService := database.New()
	db := dbService.ToGormDB()
	
	if err := db.Create(&newSchool).Error; err != nil {
		log.Printf("Error creating school: %v", err)
		http.Error(w, "Failed to create school", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newSchool); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func UpdateSchoolByUuidHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	uuidParam := chi.URLParam(r, "uuid")
	
	var updatedFields map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updatedFields); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	dbService := database.New()
	db := dbService.ToGormDB()
	
	var school model.School
	if err := db.Where("uuid = ?", uuidParam).First(&school).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "School not found", http.StatusNotFound)
		} else {
			log.Printf("Error finding school: %v", err)
			http.Error(w, "Failed to retrieve school", http.StatusInternalServerError)
		}
		return
	}
	
	if err := db.Model(&model.School{}).Where("uuid = ?", uuidParam).Updates(updatedFields).Error; err != nil {
		log.Printf("Error updating school: %v", err)
		http.Error(w, "Failed to update school", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"message": "School updated successfully"}); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
	}
}

func DeleteSchoolByUuidHandler(w http.ResponseWriter, r *http.Request) {
	uuidParam := chi.URLParam(r, "uuid")
	
	dbService := database.New()
	db := dbService.ToGormDB()
	
	var school model.School
	
	if err := db.Model(&model.School{}).Select("uuid").Where("uuid = ?", uuidParam).First(&school).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "School not found", http.StatusNotFound)
		} else {
			log.Printf("Error querying database: %v", err)
			http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		}
		return
	}
	
	if err := db.Delete(&school, "uuid = ?", uuidParam).Error; err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to delete data", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode("Data deleted successfully"); err != nil {
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
	}
}
