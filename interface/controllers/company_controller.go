package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"jobhunting-api/domain/models"
	"jobhunting-api/usecase/services"
)

type CompanyController struct {
	companyService services.CompanyService
}

func NewCompanyController(companyService services.CompanyService) *CompanyController {
	return &CompanyController{
		companyService: companyService,
	}
}

func (c *CompanyController) CreateCompany(ctx *gin.Context) {
	var req models.CompanyRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := validateCompanyRequest(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	company, err := c.companyService.CreateCompany(ctx.Request.Context(), &req)
	if err != nil {
		if isConflictError(err) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error":   "Company already exists",
				"details": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create company",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Company created successfully",
		"data":    company.ToResponse(),
	})
}

func (c *CompanyController) GetCompany(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid company ID",
		})
		return
	}

	company, err := c.companyService.GetCompany(ctx.Request.Context(), id)
	if err != nil {
		if isNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Company not found",
				"details": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get company",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": company.ToResponse(),
	})
}

func (c *CompanyController) UpdateCompany(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid company ID",
		})
		return
	}

	var req models.CompanyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := validateCompanyRequest(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	company, err := c.companyService.UpdateCompany(ctx.Request.Context(), id, &req)
	if err != nil {
		if isNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Company not found",
				"details": err.Error(),
			})
			return
		}
		if isConflictError(err) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error":   "Company name already exists",
				"details": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update company",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Company updated successfully",
		"data":    company.ToResponse(),
	})
}

func (c *CompanyController) DeleteCompany(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid company ID",
		})
		return
	}

	err = c.companyService.DeleteCompany(ctx.Request.Context(), id)
	if err != nil {
		if isNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Company not found",
				"details": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete company",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Company deleted successfully",
	})
}

func (c *CompanyController) ListCompanies(ctx *gin.Context) {
	var filter models.CompanyFilter

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	response, err := c.companyService.ListCompanies(ctx.Request.Context(), &filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list companies",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (c *CompanyController) SearchCompanies(ctx *gin.Context) {
	query := ctx.Query("q")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Search query is required",
		})
		return
	}

	companies, err := c.companyService.SearchCompanies(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to search companies",
			"details": err.Error(),
		})
		return
	}

	responses := make([]models.CompanyResponse, len(companies))
	for i, company := range companies {
		responses[i] = company.ToResponse()
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": responses,
	})
}

func validateCompanyRequest(req *models.CompanyRequest) error {
	if req.Name == "" {
		return fmt.Errorf("company name is required")
	}
	return nil
}

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "not found")
}

func isConflictError(err error) bool {
	return strings.Contains(err.Error(), "already exists")
}
