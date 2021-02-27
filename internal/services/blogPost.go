package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Bendomey/avc-server/internal/mail"
	"github.com/Bendomey/avc-server/internal/orm"
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// BlogPostService inteface holds the Tag-databse transactions of this controller
type BlogPostService interface {
	CreatePost(context context.Context, title string, details string, tag string, status string, image *string, createdBy string) (*models.BlogPost, error)
	UpdatePost(context context.Context, postID string, title *string, details *string, tag *string, status *string, image *string) (bool, error)
	DeletePost(context context.Context, postID string) (bool, error)
	ReadPost(ctx context.Context, postID string) (*models.BlogPost, error)
	ReadPosts(ctx context.Context, filterQuery *utils.FilterQuery, status *string, title *string, details *string, tag *string) ([]*models.BlogPost, error)
	ReadPostsLength(ctx context.Context, filterQuery *utils.FilterQuery, status *string, title *string, details *string, tag *string) (*int64, error)
}

// BlogPostSvc exposed the ORM to the Tag functions in the module
func BlogPostSvc(db *orm.ORM, rdb *redis.Client, mg mail.MailingService) BlogPostService {
	return &ORM{db, rdb, mg}
}

//CreatePost allow admins to create post
func (orm *ORM) CreatePost(context context.Context, title string, details string, tag string, status string, image *string, createdBy string) (*models.BlogPost, error) {
	_Post := models.BlogPost{
		Title:       title,
		Image:       image,
		Status:      status,
		TagID:       tag,
		Details:     details,
		CreatedByID: createdBy,
	}

	err := orm.DB.DB.Select("Title", "Image", "Status", "TagID", "Details", "CreatedByID").Create(&_Post).Error
	if err != nil {
		return nil, err
	}

	//later send to mail subscribers

	return &_Post, nil
}

//UpdatePost allow admins to update posts
func (orm *ORM) UpdatePost(context context.Context, postID string, title *string, details *string, tag *string, status *string, image *string) (bool, error) {
	var _Post models.BlogPost

	err := orm.DB.DB.First(&_Post, "id = ?", postID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("BlogPostNotFound")
	}

	if title != nil {
		_Post.Title = *title
	}

	if details != nil {
		_Post.Details = *details
	}

	if tag != nil {
		_Post.TagID = *tag
	}

	if status != nil {
		_Post.Status = *status
	}

	if image != nil {
		_Post.Image = image
	}

	orm.DB.DB.Save(&_Post)

	// return success
	return true, nil
}

//DeletePost allow admins to delete Posts
func (orm *ORM) DeletePost(context context.Context, postID string) (bool, error) {
	var _Post models.BlogPost

	err := orm.DB.DB.First(&_Post, "id = ?", postID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("BlogPostNotFound")
	}

	//delete
	delErr := orm.DB.DB.Delete(&_Post).Error
	if delErr != nil {
		return false, delErr
	}

	// return success
	return true, nil
}

//ReadPost calls the single in post
func (orm *ORM) ReadPost(ctx context.Context, postID string) (*models.BlogPost, error) {
	var _Post models.BlogPost

	err := orm.DB.DB.Joins("CreatedBy").First(&_Post, "blog_posts.id = ?", postID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("BlogPostNotFound")
	}

	// return success
	return &_Post, nil
}

//ReadPosts based on a query
func (orm *ORM) ReadPosts(ctx context.Context, filterQuery *utils.FilterQuery, status *string, title *string, details *string, tag *string) ([]*models.BlogPost, error) {
	var _Posts []*models.BlogPost

	_Results := orm.DB.DB

	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("blog_posts.created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if status != nil {
		_Results = _Results.Where("blog_posts.Status = ?", status)
	}

	if title != nil {
		_Results = _Results.Where("blog_posts.Title = ?", title)
	}

	if details != nil {
		_Results = _Results.Where("blog_posts.Details = ?", details)
	}

	if tag != nil {
		_Results = _Results.Where("blog_posts.TagID = ?", tag)
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("blog_posts.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("blog_posts.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("CreatedBy").
		Order(fmt.Sprintf("%s %s", filterQuery.OrderBy, filterQuery.Order)).
		Limit(filterQuery.Limit).Offset(filterQuery.Skip).
		Find(&_Posts)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return _Posts, nil
}

//ReadPostsLength retirieved the count based on a query
func (orm *ORM) ReadPostsLength(ctx context.Context, filterQuery *utils.FilterQuery, status *string, title *string, details *string, tag *string) (*int64, error) {
	var _PostsLength int64

	_Results := orm.DB.DB.Model(&models.BlogPost{})
	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("blog_posts.created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if status != nil {
		_Results = _Results.Where("blog_posts.Status = ?", status)
	}

	if title != nil {
		_Results = _Results.Where("blog_posts.Title = ?", title)
	}

	if details != nil {
		_Results = _Results.Where("blog_posts.Details = ?", details)
	}

	if tag != nil {
		_Results = _Results.Where("blog_posts.TagID = ?", tag)
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("blog_posts.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("blog_posts.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("CreatedBy").
		Count(&_PostsLength)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return &_PostsLength, nil
}
