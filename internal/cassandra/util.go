package cassandra

func mustGetEmptyStringList(listErrFunc func() ([]string, error), nonEmptyListError error) error {
	list, err := listErrFunc()
	if err != nil {
		return err
	}
	if len(list) > 0 {
		return nonEmptyListError
	}
	return nil
}