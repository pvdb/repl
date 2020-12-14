require 'test_helper'

class ReplTest < Minitest::Test
  def test_that_it_has_a_name
    refute_nil ::Repl::NAME
  end

  def test_that_it_has_a_version
    refute_nil ::Repl::VERSION
  end

  def test_that_version_matches_regexp
    assert_match Repl.v_regexp, Repl.version
  end

  def test_it_does_something_useful
    assert true # you better believe it!
  end
end
