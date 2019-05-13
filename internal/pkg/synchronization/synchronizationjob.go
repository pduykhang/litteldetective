package synchronization

// This func still not completion
//func Sync(filmID chan string, elService *film.FilmService, mongoService *film.FilmMongoService) {
//	logger := flog.New()
//	logger.SetLocal("synchronization")
//	for {
//		select {
//		case idFilm, ok := <-filmID:
//			if !ok {
//				return
//			}
//			tem, err := mongoService.GetFilmByID(idFilm)
//			if err != nil {
//
//			} else {
//				_, err := elService.InsertDataToElastic(tem)
//				if err != nil {
//
//				}
//			}
//
//		}
//	}
//}
